from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Board, Team, User, Column


class CreateBoardTests(APITestCase):
    def setUp(self):
        self.url = '/boards/'
        self.team = Team.objects.create()
        self.admin = User.objects.create(
            username='teamadmin',
            password=b'$2b$12$lrkDnrwXSBU.YJvdzbpAWOd9GhwHJGVYafRXTHct2gm3akPJ'
                     b'gB5Zq',
            is_admin=True,
            team=self.team
        )
        self.member = User.objects.create(
            username='teammember',
            password=b'$2b$12$RonFQ1/18JiCN8yFeBaxKOsVbxLdcehlZ4e0r9gtZbARqEVU'
                     b'HHEoK',
            is_admin=False,
            team=self.team
        )
        self.admin_token = '$2b$12$TVdxI.a.ZlOkhH1/mZQ/IOHmKxklQJWiB0n6ZSg2R' \
                           'JJO17thjLOdy'
        self.member_token = '$2b$12$xnIJLzpgNV12O80XsakMjezCFqwIphdBy5ziJ9Eb' \
                            '9stnDZze19Ude'

    def test_success(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url,
                                    {'team_id': self.team.id},
                                    HTTP_AUTH_USER=self.admin.username,
                                    HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.data.get('msg'),
                         'Board creation successful.')
        board = Board.objects.get(id=response.data.get('board_id'))
        columns = Column.objects.filter(board=board.id)
        self.assertEqual(len(columns), 4)
        self.assertEqual(Board.objects.count(), initial_count + 1)

    def test_username_blank(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url,
                                    {'team_id': self.team.id},
                                    HTTP_AUTH_USER='',
                                    HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': ErrorDetail(string='Username cannot be empty.',
                                    code='blank')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_username_invalid(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url,
                                    {'team_id': self.team.id},
                                    HTTP_AUTH_USER='invalidio',
                                    HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': ErrorDetail(string='Invalid username.', code='invalid')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_user_not_admin(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url,
                                    {'team_id': self.team.id},
                                    HTTP_AUTH_USER=self.member.username,
                                    HTTP_AUTH_TOKEN=self.member_token)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': ErrorDetail(
                string='Only the team admin can create a board.',
                code='not_authorized'
            )
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_team_id_blank(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url,
                                    {'team_id': ''},
                                    HTTP_AUTH_USER=self.admin.username,
                                    HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team ID cannot be empty.',
                                   code='blank')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_team_not_found(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url,
                                    {'team_id': '123'},
                                    HTTP_AUTH_USER=self.admin.username,
                                    HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team not found.', code='not_found')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_token_empty(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url,
                                    {'team_id': self.team.id},
                                    HTTP_AUTH_USER=self.admin.username,
                                    HTTP_AUTH_TOKEN='')
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'token': ErrorDetail(
                string='Authentication token cannot be empty.',
                code='blank',
            )
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_token_invalid(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url,
                                    {'username': self.admin.username,
                                     'token': 'ASDKFJ!FJ_012rjpiwajfosi12311@',
                                     'team_id': self.team.id})
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'token': ErrorDetail(
                string='Invalid authentication token.',
                code='invalid',
            )
        })
        self.assertEqual(Board.objects.count(), initial_count)