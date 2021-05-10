from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
import status

from ..models import Board, Team
from ..util import create_admin, create_member
from ..validation.auth import \
    not_authenticated_response, authorization_error


class PatchBoardTests(APITestCase):
    endpoint = '/boards/?id='

    def setUp(self):
        team = Team.objects.create()
        self.admin = create_admin(team)
        self.member = create_member(team)
        self.board = Board.objects.create(name='Some Board', team=team)
        self.wrong_admin = create_admin(Team.objects.create(), '1')

    def test_success(self):
        response = self.client.patch(f'{self.endpoint}{self.board.id}',
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'],
                                     format='json')
        print(f'successresponse: {response.data}')
        self.assertEqual(response.status_code, status.HTTP_200_OK)
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
                                     HTTP_AUTH_TOKEN=self.admin['token'],
                                     format='json')
        self.assertEqual(response.status_code, status.HTTP_400_BAD_REQUEST)
        self.assertEqual(response.data, {
            'id': [ErrorDetail(string='Board ID cannot be null.',
                               code='null')]
        })
        self.assertEqual(Board.objects.get(id=self.board.id).name,
                         'Some Board')

    def test_board_id_invalid(self):
        response = self.client.patch(f'{self.endpoint}sadfj',
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'],
                                     format='json')
        self.assertEqual(response.status_code, status.HTTP_400_BAD_REQUEST)
        self.assertEqual(response.data, {
            'id': [ErrorDetail(string='Board ID must be a number.',
                               code='invalid')]
        })
        self.assertEqual(Board.objects.get(id=self.board.id).name,
                         'Some Board')

    def test_board_not_found(self):
        response = self.client.patch(f'{self.endpoint}1231231',
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'],
                                     format='json')
        self.assertEqual(response.status_code, status.HTTP_404_NOT_FOUND)
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
                                     HTTP_AUTH_TOKEN=self.admin['token'],
                                     format='json')
        self.assertEqual(response.status_code, status.HTTP_400_BAD_REQUEST)
        self.assertEqual(response.data, {
            'name': [ErrorDetail(string='Board name cannot be blank.',
                                 code='blank')]
        })
        self.assertEqual(Board.objects.get(id=self.board.id).name,
                         'Some Board')

    def test_auth_user_empty(self):
        response = self.client.patch(f'{self.endpoint}{self.board.id}',
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER='',
                                     HTTP_AUTH_TOKEN=self.admin['token'],
                                     format='json')
        self.assertEqual(response.status_code,
                         not_authenticated_response.status_code)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_user_invalid(self):
        response = self.client.patch(f'{self.endpoint}{self.board.id}',
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER='invalidusername',
                                     HTTP_AUTH_TOKEN=self.admin['token'],
                                     format='json')
        self.assertEqual(response.status_code,
                         not_authenticated_response.status_code)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_token_empty(self):
        response = self.client.patch(f'{self.endpoint}{self.board.id}',
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN='',
                                     format='json')
        self.assertEqual(response.status_code,
                         not_authenticated_response.status_code)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_token_invalid(self):
        response = self.client.patch(f'{self.endpoint}{self.board.id}',
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN='ASDKFJ!FJ_012rjpajfosia',
                                     format='json')
        self.assertEqual(response.status_code,
                         not_authenticated_response.status_code)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_wrong_team(self):
        initial_count = Board.objects.count()
        response = self.client.patch(
            f'{self.endpoint}{self.board.id}',
            {'name': 'New Title'},
            HTTP_AUTH_USER=self.wrong_admin['username'],
            HTTP_AUTH_TOKEN=self.wrong_admin['token'],
            format='json'
        )
        self.assertEqual(response.status_code, status.HTTP_401_UNAUTHORIZED)
        self.assertEqual(response.data, authorization_error.detail)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_unauthorized(self):
        response = self.client.patch(f'{self.endpoint}{self.board.id}',
                                     {'name': 'New Title'},
                                     HTTP_AUTH_USER=self.member['username'],
                                     HTTP_AUTH_TOKEN=self.member['token'],
                                     format='json')
        self.assertEqual(response.status_code, status.HTTP_401_UNAUTHORIZED)
        self.assertEqual(response.data, authorization_error.detail)
