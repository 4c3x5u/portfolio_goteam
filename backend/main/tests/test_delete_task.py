from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Team, Board, Column, Task
from ..helpers import UserHelper
from ..validation.val_auth import authentication_error, authorization_error


class DeleteTaskTests(APITestCase):
    endpoint = '/tasks/?id='

    def setUp(self):
        team = Team.objects.create()
        board = Board.objects.create(team=team)
        column = Column.objects.create(order=0, board=board)
        self.task = Task.objects.create(title='Do Something!',
                                        order=0,
                                        column=column)

        user_helper = UserHelper(team)
        self.member = user_helper.create()
        self.admin = user_helper.create(is_admin=True)

        wrong_user_helper = UserHelper(Team.objects.create())
        self.wrong_admin = wrong_user_helper.create(is_admin=True)

    def test_success(self):
        initial_count = Task.objects.count()
        response = self.client.delete(f'{self.endpoint}{self.task.id}',
                                      HTTP_AUTH_USER=self.admin['username'],
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': 'Task deleted successfully.',
            'id': self.task.id
        })
        self.assertEqual(Task.objects.count(), initial_count - 1)

    def test_task_id_blank(self):
        initial_count = Task.objects.count()
        response = self.client.delete(self.endpoint,
                                      HTTP_AUTH_USER=self.admin['username'],
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'task': [ErrorDetail(string='Task ID cannot be null.',
                                 code='null')]
        })
        self.assertEqual(Task.objects.count(), initial_count)

    def test_task_id_invalid(self):
        initial_count = Task.objects.count()
        response = self.client.delete(f'{self.endpoint}qwerty',
                                      HTTP_AUTH_USER=self.admin['username'],
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        print(f'id invalid response data {response.data}')
        self.assertEqual(response.data, {
            'task': [ErrorDetail(string='Task ID must be a number.',
                                 code='incorrect_type')]
        })
        self.assertEqual(Task.objects.count(), initial_count)

    def test_task_not_found(self):
        initial_count = Task.objects.count()
        response = self.client.delete(f'{self.endpoint}123141',
                                      HTTP_AUTH_USER=self.admin['username'],
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'task': [ErrorDetail(string='Task does not exist.',
                                 code='does_not_exist')]
        })
        self.assertEqual(Task.objects.count(), initial_count)

    def test_auth_token_empty(self):
        response = self.client.delete(f'{self.endpoint}{self.task.id}',
                                      HTTP_AUTH_USER=self.admin['username'],
                                      HTTP_AUTH_TOKEN='')
        self.assertEqual(response.status_code,
                         authentication_error.status_code)
        self.assertEqual(response.data, authentication_error.detail)

    def test_auth_token_invalid(self):
        response = self.client.delete(
            f'{self.endpoint}{self.task.id}',
            HTTP_AUTH_USER=self.admin['username'],
            HTTP_AUTH_TOKEN='ASDKFJ!FJ_012rjpiwajfos'
        )
        self.assertEqual(response.status_code,
                         authentication_error.status_code)
        self.assertEqual(response.data, authentication_error.detail)

    def test_auth_user_blank(self):
        response = self.client.delete(f'{self.endpoint}{self.task.id}',
                                      HTTP_AUTH_USER='',
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code,
                         authentication_error.status_code)
        self.assertEqual(response.data, authentication_error.detail)

    def test_auth_user_invalid(self):
        response = self.client.delete(f'{self.endpoint}{self.task.id}',
                                      HTTP_AUTH_USER='invalidio',
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code,
                         authentication_error.status_code)
        self.assertEqual(response.data, authentication_error.detail)

    def test_wrong_team(self):
        initial_count = Board.objects.count()
        response = self.client.delete(
            f'{self.endpoint}{self.task.id}',
            HTTP_AUTH_USER=self.wrong_admin['username'],
            HTTP_AUTH_TOKEN=self.wrong_admin['token']
        )
        self.assertEqual(response.status_code,
                         authorization_error.status_code)
        self.assertEqual(response.data, authorization_error.detail)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_unauthorized(self):
        response = self.client.delete(f'{self.endpoint}{self.task.id}',
                                      HTTP_AUTH_USER=self.member['username'],
                                      HTTP_AUTH_TOKEN=self.member['token'])
        self.assertEqual(response.status_code, authorization_error.status_code)
        self.assertEqual(response.data, authorization_error.detail)

