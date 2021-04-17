from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Board, Team
from ..util import new_admin


class PatchBoardTests(APITestCase):
    endpoint = '/boards/?id='

    def setUp(self):
        team = Team.objects.create()
        self.admin = new_admin(team)
        self.board = Board.objects.create(name='Some Board',
                                          team=team)

    def test_success(self):
        response = self.client.patch(f'{self.endpoint}{self.board.id}',
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': 'Board updated successfuly.',
            'id': self.board.id,
        })
        self.assertEqual(Board.objects.get(id=self.board.id).name,
                         'New Title')

    def test_board_id_empty(self):
        response = self.client.patch(self.endpoint,
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'board_id': ErrorDetail(string='Board ID cannot be empty.',
                                    code='blank')
        })
        self.assertEqual(Board.objects.get(id=self.board.id).name,
                         'Some Board')

    def test_board_id_invalid(self):
        response = self.client.patch(f'{self.endpoint}sadfj',
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'board_id': ErrorDetail(string='Board ID must be a number.',
                                    code='invalid')
        })
        self.assertEqual(Board.objects.get(id=self.board.id).name,
                         'Some Board')

    def test_board_not_found(self):
        response = self.client.patch(f'{self.endpoint}1231231',
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, {
            'board_id': ErrorDetail(string='Board not found.',
                                    code='not_found')
        })
        self.assertEqual(Board.objects.get(id=self.board.id).name,
                         'Some Board')
