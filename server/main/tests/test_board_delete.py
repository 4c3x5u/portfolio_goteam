from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Team, Board
from server.main.helpers.user_helper import UserHelper
from server.main.helpers.auth_helper import AuthHelper


class DeleteBoardTests(APITestCase):
    endpoint = '/boards/?id='

    def setUp(self):
        team = Team.objects.create()
        boards = Board.objects.bulk_create([Board(team=team)
                                            for _ in range(0, 4)])
        self.board = boards[0]
        user_helper = UserHelper(team)
        self.member = user_helper.create_user()
        self.admin = user_helper.create_user(is_admin=True)

        wrong_user_helper = UserHelper(Team.objects.create())
        self.wrong_admin = wrong_user_helper.create_user(is_admin=True)

    def test_success(self):
        initial_count = Board.objects.count()
        response = self.client.delete(f'{self.endpoint}{self.board.id}',
                                      HTTP_AUTH_USER=self.admin['username'],
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        print(f'success: {response.data}')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': 'Board deleted successfully.',
            'id': self.board.id,
        })
        self.assertEqual(Board.objects.count(), initial_count - 1)

    def test_board_id_blank(self):
        initial_count = Board.objects.count()
        response = self.client.delete(self.endpoint,
                                      HTTP_AUTH_USER=self.admin['username'],
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'board': [ErrorDetail(string='Board ID cannot be null.',
                                  code='null')]
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_board_id_invalid(self):
        initial_count = Board.objects.count()
        response = self.client.delete(f'{self.endpoint}qwerty',
                                      HTTP_AUTH_USER=self.admin['username'],
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'board': [ErrorDetail(string='Board ID must be a number.',
                                  code='incorrect_type')]
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_board_not_found(self):
        initial_count = Board.objects.count()
        response = self.client.delete(f'{self.endpoint}123141',
                                      HTTP_AUTH_USER=self.admin['username'],
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'board': [ErrorDetail(string='Board does not exist.',
                                  code='does_not_exist')]
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_token_empty(self):
        initial_count = Board.objects.count()
        response = self.client.delete(f'{self.endpoint}{self.board.id}',
                                      HTTP_AUTH_USER=self.admin['username'],
                                      HTTP_AUTH_TOKEN='')
        self.assertEqual(response.status_code,
                         AuthHelper.AUTHENTICATION_ERROR.status_code)
        self.assertEqual(response.data,
                         AuthHelper.AUTHENTICATION_ERROR.detail)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_token_invalid(self):
        initial_count = Board.objects.count()
        response = self.client.delete(f'{self.endpoint}{self.board.id}',
                                      HTTP_AUTH_USER=self.admin['username'],
                                      HTTP_AUTH_TOKEN='ASDKFJ!FJ_012rjpiwajfos')
        self.assertEqual(response.status_code,
                         AuthHelper.AUTHENTICATION_ERROR.status_code)
        self.assertEqual(response.data,
                         AuthHelper.AUTHENTICATION_ERROR.detail)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_user_blank(self):
        initial_count = Board.objects.count()
        response = self.client.delete(f'{self.endpoint}{self.board.id}',
                                      HTTP_AUTH_USER='',
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code,
                         AuthHelper.AUTHENTICATION_ERROR.status_code)
        self.assertEqual(response.data,
                         AuthHelper.AUTHENTICATION_ERROR.detail)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_user_invalid(self):
        initial_count = Board.objects.count()
        response = self.client.delete(f'{self.endpoint}{self.board.id}',
                                      HTTP_AUTH_USER='invalidio',
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code,
                         AuthHelper.AUTHENTICATION_ERROR.status_code)
        self.assertEqual(response.data,
                         AuthHelper.AUTHENTICATION_ERROR.detail)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_wrong_team(self):
        initial_count = Board.objects.count()
        response = self.client.delete(
            f'{self.endpoint}{self.board.id}',
            HTTP_AUTH_USER=self.wrong_admin['username'],
            HTTP_AUTH_TOKEN=self.wrong_admin['token'],
        )
        self.assertEqual(response.status_code,
                         AuthHelper.AUTHORIZATION_ERROR.status_code)
        self.assertEqual(response.data,
                         AuthHelper.AUTHORIZATION_ERROR.detail)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_unauthorized(self):
        initial_count = Board.objects.count()
        response = self.client.delete(f'{self.endpoint}{self.board.id}',
                                      HTTP_AUTH_USER=self.member['username'],
                                      HTTP_AUTH_TOKEN=self.member['token'])
        self.assertEqual(response.status_code,
                         AuthHelper.AUTHORIZATION_ERROR.status_code)
        self.assertEqual(response.data,
                         AuthHelper.AUTHORIZATION_ERROR.detail)
        self.assertEqual(Board.objects.count(), initial_count)
