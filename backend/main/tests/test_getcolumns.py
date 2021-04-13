from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Board, Team, Column, Task
from ..util import new_admin


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
