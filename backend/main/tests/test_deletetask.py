from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Team, Board, Column, Task
from ..util import new_admin, new_member, not_authenticated_response


class DeleteTaskTests(APITestCase):
    endpoint = '/tasks/?id='

    def setUp(self):
        team = Team.objects.create()
        board = Board.objects.create(team=team)
        column = Column.objects.create(order=0, board=board)
        self.task = Task.objects.create(title='Do Something!',
                                        order=0,
                                        column=column)
        self.admin = new_admin(team)
        self.member = new_member(team)

    def test_success(self):
        initial_count = Task.objects.count()
        response = self.client.delete(f'{self.endpoint}{self.task.id}',
                                      HTTP_AUTH_USER=self.admin['username'],
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': 'Task deleted successfully.',
            'id': str(self.task.id),
        })
        self.assertEqual(Task.objects.count(), initial_count - 1)

    def test_task_id_blank(self):
        initial_count = Task.objects.count()
        response = self.client.delete(self.endpoint,
                                      HTTP_AUTH_USER=self.admin['username'],
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'task_id': ErrorDetail(string='Task ID cannot be empty.',
                                   code='blank')
        })
        self.assertEqual(Task.objects.count(), initial_count)

    def test_task_id_invalid(self):
        initial_count = Task.objects.count()
        response = self.client.delete(f'{self.endpoint}qwerty',
                                      HTTP_AUTH_USER=self.admin['username'],
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'task_id': ErrorDetail(string='Task ID must be a number.',
                                   code='invalid')
        })
        self.assertEqual(Task.objects.count(), initial_count)

    def test_task_not_found(self):
        initial_count = Task.objects.count()
        response = self.client.delete(f'{self.endpoint}123141',
                                      HTTP_AUTH_USER=self.admin['username'],
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, {
            'task_id': ErrorDetail(string='Task not found.',
                                   code='not_found')
        })
        self.assertEqual(Task.objects.count(), initial_count)


    def test_auth_token_empty(self):
        response = self.client.delete(f'{self.endpoint}{self.task.id}',
                                      HTTP_AUTH_USER=self.admin['username'],
                                      HTTP_AUTH_TOKEN='')
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_token_invalid(self):
        response = self.client.delete(f'{self.endpoint}{self.task.id}',
                                      HTTP_AUTH_USER=self.admin['username'],
                                      HTTP_AUTH_TOKEN='ASDKFJ!FJ_012rjpiwajfos')
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_user_blank(self):
        response = self.client.delete(f'{self.endpoint}{self.task.id}',
                                      HTTP_AUTH_USER='',
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_user_invalid(self):
        response = self.client.delete(f'{self.endpoint}{self.task.id}',
                                      HTTP_AUTH_USER='invalidio',
                                      HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_unauthorized(self):
        response = self.client.delete(f'{self.endpoint}{self.task.id}',
                                      HTTP_AUTH_USER=self.member['username'],
                                      HTTP_AUTH_TOKEN=self.member['token'])
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, {
            'auth': ErrorDetail(string='The user is not an admin.',
                                code='not_authorized')
        })
