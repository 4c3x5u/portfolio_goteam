from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Board, Team, User


class ListBoardsTests(APITestCase):
    def setUp(self):
        self.base_url = '/boards/?team_id='
        self.team = Team.objects.create()
        self.boards = [
            Board.objects.create(team_id=self.team.id) for _ in range(0, 3)
        ]
        self.team_id = str(self.team.id)
        self.member = User.objects.create(
            username='teammember',
            password=b'$2b$12$RonFQ1/18JiCN8yFeBaxKOsVbxLdcehlZ4e0r9gtZbARqEVU'
                     b'HHEoK',
            is_admin=False,
            team=self.team
        )
        self.member_token = '$2b$12$xnIJLzpgNV12O80XsakMjezCFqwIphdBy5ziJ9Eb' \
                            '9stnDZze19Ude'
        self.forbidden_response = {
            'auth': ErrorDetail(string="Authentication failure.",
                                code='not_authenticated')
        }

    def test_success(self):
        initial_count = Board.objects.count()
        response = self.client.get(self.base_url + self.team_id,
                                   HTTP_AUTH_USER=self.member.username,
                                   HTTP_AUTH_TOKEN=self.member_token)
        self.assertEqual(response.status_code, 200)
        boards = response.data.get('boards')
        self.assertTrue(boards)
        self.assertTrue(boards.count, 3)
        for board in boards:
            self.assertEqual(board.get('team'), self.team.id)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_boards_not_found_member(self):
        initial_count = Board.objects.count()
        team = Team.objects.create()
        response = self.client.get(self.base_url + str(team.id),
                                   HTTP_AUTH_USER=self.member.username,
                                   HTTP_AUTH_TOKEN=self.member_token)
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Boards not found.',
                                   code='not_found')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    # if no boards are found for a team admin,
    # a new one is created and returned
    def test_boards_not_found_admin(self):
        initial_count = Board.objects.count()
        team = Team.objects.create()
        user = User.objects.create(
            username='teamadmin',
            password=b'$2b$12$lrkDnrwXSBU.YJvdzbpAWOd9GhwHJGVYafRXTHct2gm3akPJ'
                     b'gB5Zq',
            is_admin=True,
            team=self.team
        )
        token = '$2b$12$TVdxI.a.ZlOkhH1/mZQ/IOHmKxklQJWiB0n6ZSg2R' \
                           'JJO17thjLOdy'
        response = self.client.get(self.base_url + str(team.id),
                                   HTTP_AUTH_USER=user.username,
                                   HTTP_AUTH_TOKEN=token)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(len(response.data.get('boards')), 1)
        self.assertEqual(Board.objects.count(), initial_count + 1)

    def test_team_id_empty(self):
        initial_count = Board.objects.count()
        response = self.client.get(self.base_url,
                                   HTTP_AUTH_USER=self.member.username,
                                   HTTP_AUTH_TOKEN=self.member_token)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team ID cannot be empty.',
                                   code='blank')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_team_not_found(self):
        initial_count = Board.objects.count()
        response = self.client.get(self.base_url + '123',
                                   HTTP_AUTH_USER=self.member.username,
                                   HTTP_AUTH_TOKEN=self.member_token)
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team not found.', code='not_found')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    # TODO:
    #  [X]. new board created for admin
    #  [X]. auth user empty
    #  [X]. auth user invalid
    #  [ ]. auth token empty
    #  [ ]. auth token invalid

    def test_auth_user_empty(self):
        initial_count = Board.objects.count()
        response = self.client.get(self.base_url + self.team_id,
                                   HTTP_AUTH_USER='',
                                   HTTP_AUTH_TOKEN=self.member_token)
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, self.forbidden_response)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_user_invalid(self):
        initial_count = Board.objects.count()
        response = self.client.get(self.base_url + self.team_id,
                                   HTTP_AUTH_USER='invalidusername',
                                   HTTP_AUTH_TOKEN=self.member_token)
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, self.forbidden_response)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_token_empty(self):
        initial_count = Board.objects.count()
        response = self.client.get(self.base_url + self.team_id,
                                   HTTP_AUTH_USER=self.member.username,
                                   HTTP_AUTH_TOKEN='')
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, self.forbidden_response)
        self.assertEqual(Board.objects.count(), initial_count)
