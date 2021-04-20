from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Board, Team
from ..util import new_admin, new_member, not_authenticated_response


class PatchBoardTests(APITestCase):
    endpoint = '/boards/?id='

    def setUp(self):
        team = Team.objects.create()
        self.admin = new_admin(team)
        self.member = new_member(team)
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

    def test_board_name_blank(self):
        response = self.client.patch(f'{self.endpoint}{self.board.id}',
                                     {'name': ''},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'name': [ErrorDetail(string='Board name cannot be empty.',
                                 code='blank')]
        })
        self.assertEqual(Board.objects.get(id=self.board.id).name,
                         'Some Board')

    def test_auth_user_empty(self):
        response = self.client.patch(f'{self.endpoint}{self.board.id}',
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER='',
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code,
                         not_authenticated_response.status_code)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_user_invalid(self):
        response = self.client.patch(f'{self.endpoint}{self.board.id}',
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER='invalidusername',
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code,
                         not_authenticated_response.status_code)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_token_empty(self):
        response = self.client.patch(f'{self.endpoint}{self.board.id}',
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN='')
        self.assertEqual(response.status_code,
                         not_authenticated_response.status_code)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_token_invalid(self):
        response = self.client.patch(f'{self.endpoint}{self.board.id}',
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN='ASDKFJ!FJ_012rjpajfosia')
        self.assertEqual(response.status_code,
                         not_authenticated_response.status_code)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_unauthorized(self):
        response = self.client.patch(f'{self.endpoint}{self.board.id}',
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER=self.member['username'],
                                     HTTP_AUTH_TOKEN=self.member['token'])
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, {
            'auth': ErrorDetail(string='You must be an admin to do this.',
                                code='not_authorized')
        })
