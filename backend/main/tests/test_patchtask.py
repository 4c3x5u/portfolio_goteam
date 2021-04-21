from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Task, Column, Board, Team
from ..util import (
    new_member, new_admin, not_authenticated_response, not_authorized_response
)


class UpdateTaskTests(APITestCase):
    endpoint = '/tasks/?id='

    def setUp(self):
        team = Team.objects.create()
        self.member = new_member(team)
        self.admin = new_admin(team)
        self.task = Task.objects.create(
            title="Task Title",
            order=0,
            column=Column.objects.create(
                order=0,
                board=Board.objects.create(
                    team=team
                )
            )
        )
        self.wrong_admin = new_admin(Team.objects.create(), '1')

    def help_test_success(self, task_id, request_data):
        response = self.client.patch(f'{self.endpoint}{task_id}',
                                     request_data,
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {'msg': 'Task update successful.',
                                         'id': self.task.id})
        self.assertEqual(self.task.id, response.data.get('id'))

    def test_title_success(self):
        request_data = {'title': 'New Title'}
        self.help_test_success(self.task.id, request_data)
        self.assertEqual(Task.objects.get(id=self.task.id).title,
                         request_data.get('title'))

    def test_title_blank(self):
        response = self.client.patch(f'{self.endpoint}{self.task.id}',
                                     {'title': ''},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'title': ErrorDetail(string='Task title cannot be empty.',
                                 code='blank')
        })
        self.assertEqual(Task.objects.get(id=self.task.id).title,
                         self.task.title)

    def test_order_success(self):
        request_data = {'order': 10}
        self.help_test_success(self.task.id, request_data)
        self.assertEqual(Task.objects.get(id=self.task.id).order,
                         request_data.get('order'))

    def test_order_blank(self):
        response = self.client.patch(f'{self.endpoint}{self.task.id}',
                                     {'order': ''},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'order': ErrorDetail(string='Task order cannot be empty.',
                                 code='blank')
        })
        self.assertEqual(Task.objects.get(id=self.task.id).order,
                         self.task.order)

    def test_column_success(self):
        another_column = Column.objects.create(
            order=0,
            board=Board.objects.create(
                team=Team.objects.create()
            )
        )
        request_data = {'column': another_column.id}
        self.help_test_success(self.task.id, request_data)
        self.assertEqual(Task.objects.get(id=self.task.id).column.id,
                         request_data.get('column'))

    def test_column_blank(self):
        response = self.client.patch(f'{self.endpoint}{self.task.id}',
                                     {'column': ''},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'column_id': ErrorDetail(string='Column ID cannot be empty.',
                                     code='blank')
        })
        self.assertEqual(Task.objects.get(id=self.task.id).column,
                         self.task.column)

    def test_column_not_found(self):
        response = self.client.patch(f'{self.endpoint}{self.task.id}',
                                     {'column': '123123'},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, {
            'column_id': ErrorDetail(string='Column not found.',
                                     code='not_found')
        })
        self.assertEqual(Task.objects.get(id=self.task.id).column,
                         self.task.column)

    def test_auth_token_empty(self):
        initial_count = Task.objects.count()
        response = self.client.patch(f'{self.endpoint}{self.task.id}',
                                     {'title': 'New Title'},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN='')
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_token_invalid(self):
        initial_count = Task.objects.count()
        response = self.client.patch(f'{self.endpoint}{self.task.id}',
                                     {'title': 'New Title'},
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN='ASDKFJ!FJ_012rjpiwajfosi')
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_user_blank(self):
        initial_count = Task.objects.count()
        response = self.client.patch(f'{self.endpoint}{self.task.id}',
                                     {'title': 'New Title'},
                                     HTTP_AUTH_USER='',
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_user_invalid(self):
        initial_count = Task.objects.count()
        response = self.client.patch(f'{self.endpoint}{self.task.id}',
                                     {'title': 'New Title'},
                                     HTTP_AUTH_USER='invalidio',
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_not_admin(self):
        initial_count = Task.objects.count()
        response = self.client.patch(f'{self.endpoint}{self.task.id}',
                                     {'title': 'New Title'},
                                     HTTP_AUTH_USER=self.member['username'],
                                     HTTP_AUTH_TOKEN=self.member['token'])
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authorized_response.data)
        self.assertEqual(Task.objects.count(), initial_count)

    def test_wrong_team(self):
        initial_count = Task.objects.count()
        response = self.client.patch(
            f'{self.endpoint}{self.task.id}',
            {'title': 'New Title'},
            HTTP_AUTH_USER=self.wrong_admin['username'],
            HTTP_AUTH_TOKEN=self.wrong_admin['token']
        )
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)
        self.assertEqual(Task.objects.count(), initial_count)
