from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Board, Team, Column
from ..util import new_admin, not_authenticated_response


class GetColumnsTests(APITestCase):
    endpoint = '/columns/?board_id='

    def setUp(self):
        self.team = Team.objects.create()
        self.admin = new_admin(self.team)
        self.board = Board.objects.create(team_id=self.team.id)
        self.columns = [
            Column.objects.create(
                order=i, board=self.board
            ) for i in range(0, 4)
        ]

    def test_success(self):
        response = self.client.get(f'{self.endpoint}{self.board.id}',
                                   HTTP_AUTH_USER=self.admin['username'],
                                   HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 200)
        columns = response.data.get('columns')
        self.assertTrue(columns)
        self.assertTrue(columns.count, 4)
        for i in range(0, 4):
            self.assertEqual(self.columns[i].id, columns[i].get('id'))

    def test_board_id_empty(self):
        response = self.client.get(self.endpoint,
                                   HTTP_AUTH_USER=self.admin['username'],
                                   HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'board_id': ErrorDetail(string='Board ID cannot be empty.',
                                    code='blank')
        })

    def test_board_id_invalid(self):
        response = self.client.get(f'{self.endpoint}asdf',
                                   HTTP_AUTH_USER=self.admin['username'],
                                   HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'board_id': ErrorDetail(string='Board ID must be a number.',
                                    code='invalid')
        })

    def test_board_not_found(self):
        response = self.client.get(f'{self.endpoint}123123',
                                   HTTP_AUTH_USER=self.admin['username'],
                                   HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, {
            'board_id': ErrorDetail(string='Board not found.', code='not_found')
        })

    def test_auth_user_empty(self):
        response = self.client.get(f'{self.endpoint}{self.board.id}',
                                   HTTP_AUTH_USER='',
                                   HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code,
                         not_authenticated_response.status_code)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_user_invalid(self):
        response = self.client.get(f'{self.endpoint}{self.board.id}',
                                   HTTP_AUTH_USER='invalidusername',
                                   HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code,
                         not_authenticated_response.status_code)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_token_empty(self):
        response = self.client.get(f'{self.endpoint}{self.board.id}',
                                   HTTP_AUTH_USER=self.admin['username'],
                                   HTTP_AUTH_TOKEN='')
        self.assertEqual(response.status_code,
                         not_authenticated_response.status_code)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_token_invalid(self):
        response = self.client.get(f'{self.endpoint}{self.board.id}',
                                   HTTP_AUTH_USER=self.admin['username'],
                                   HTTP_AUTH_TOKEN='ASDKFJ!FJ_012rjpiwajfosia')
        self.assertEqual(response.status_code,
                         not_authenticated_response.status_code)
        self.assertEqual(response.data, not_authenticated_response.data)
