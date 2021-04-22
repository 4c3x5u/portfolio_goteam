from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Board, Team, Column
from ..util import create_member, create_admin
from ..validation.val_auth import not_authenticated_response


class GetBoardsTests(APITestCase):
    endpoint = '/boards/?team_id='

    def setUp(self):
        self.team = Team.objects.create()
        self.member = create_member(self.team)
        self.boards = [
            Board.objects.create(team_id=self.team.id) for _ in range(0, 3)
        ]
        self.wrong_member = create_member(Team.objects.create(), '1')

    def test_success(self):
        initial_count = Board.objects.count()
        response = self.client.get(f'{self.endpoint}{self.team.id}',
                                   HTTP_AUTH_USER=self.member['username'],
                                   HTTP_AUTH_TOKEN=self.member['token'])
        self.assertEqual(response.status_code, 200)
        self.assertTrue(len(response.data), 3)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_boards_not_found_member(self):
        initial_count = Board.objects.count()
        team = Team.objects.create()
        member = create_member(team, '2')
        response = self.client.get(f'{self.endpoint}{team.id}',
                                   HTTP_AUTH_USER=member['username'],
                                   HTTP_AUTH_TOKEN=member['token'])
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Boards not found.',
                                   code='not_found')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    # if no boards are found for a team admin,
    # a new one is created and returned
    def test_boards_not_found_admin(self):
        initial_board_count = Board.objects.count()
        initial_columns_count = Column.objects.count()
        team = Team.objects.create()
        admin = create_admin(team)
        response = self.client.get(self.endpoint + str(team.id),
                                   HTTP_AUTH_USER=admin['username'],
                                   HTTP_AUTH_TOKEN=admin['token'])
        print(response.data)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(len(response.data), 1)
        self.assertEqual(response.data[0].get('name'), 'New Board')
        self.assertEqual(Board.objects.count(), initial_board_count + 1)
        self.assertEqual(Column.objects.count(), initial_columns_count + 4)

    def test_team_id_empty(self):
        initial_count = Board.objects.count()
        response = self.client.get(self.endpoint,
                                   HTTP_AUTH_USER=self.member['username'],
                                   HTTP_AUTH_TOKEN=self.member['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team ID cannot be empty.',
                                   code='blank')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_team_not_found(self):
        initial_count = Board.objects.count()
        response = self.client.get(self.endpoint + '123',
                                   HTTP_AUTH_USER=self.member['username'],
                                   HTTP_AUTH_TOKEN=self.member['token'])
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team not found.', code='not_found')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_user_empty(self):
        initial_count = Board.objects.count()
        response = self.client.get(f'{self.endpoint}{self.team.id}',
                                   HTTP_AUTH_USER='',
                                   HTTP_AUTH_TOKEN=self.member['token'])
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_user_invalid(self):
        initial_count = Board.objects.count()
        response = self.client.get(f'{self.endpoint}{self.team.id}',
                                   HTTP_AUTH_USER='invalidusername',
                                   HTTP_AUTH_TOKEN=self.member['token'])
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_token_empty(self):
        initial_count = Board.objects.count()
        response = self.client.get(f'{self.endpoint}{self.team.id}',
                                   HTTP_AUTH_USER=self.member['username'],
                                   HTTP_AUTH_TOKEN='')
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_token_invalid(self):
        initial_count = Board.objects.count()
        response = self.client.get(f'{self.endpoint}{self.team.id}',
                                   HTTP_AUTH_USER=self.member['username'],
                                   HTTP_AUTH_TOKEN='ASDKFJ!FJ_012rjpiwajfosia')
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_wrong_team(self):
        initial_count = Board.objects.count()
        response = self.client.get(
            f'{self.endpoint}{self.team.id}',
            HTTP_AUTH_USER=self.wrong_member['username'],
            HTTP_AUTH_TOKEN=self.wrong_member['token'],
        )
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)
        self.assertEqual(Board.objects.count(), initial_count)
