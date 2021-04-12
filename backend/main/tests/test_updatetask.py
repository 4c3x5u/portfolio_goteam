from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Task, Column, Board, Team, User

# TODO:
#   [X] Make existing tests pass
#   [X] Add authentication tests
#   [ ] Add an authorization test


class UpdateTaskTests(APITestCase):
    def setUp(self):
        self.url = '/tasks/'
        team = Team.objects.create()
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
        self.admin = User.objects.create(
            username='teamadmin',
            password=b'$2b$12$lrkDnrwXSBU.YJvdzbpAWOd9GhwHJGVYafRXTHct2gm3akPJ'
                     b'gB5Zq',
            is_admin=True,
            team=team
        )
        self.member = User.objects.create(
            username='teammember',
            password=b'$2b$12$RonFQ1/18JiCN8yFeBaxKOsVbxLdcehlZ4e0r9gtZbARqEVU'
                     b'HHEoK',
            is_admin=False,
            team=team
        )
        self.admin_token = '$2b$12$TVdxI.a.ZlOkhH1/mZQ/IOHmKxklQJWiB0n6ZSg2R' \
                           'JJO17thjLOdy'
        self.member_token = '$2b$12$xnIJLzpgNV12O80XsakMjezCFqwIphdBy5ziJ9Eb' \
                            '9stnDZze19Ude'
        self.forbidden_response = {
            'auth': ErrorDetail(string="Authentication failure.",
                                code='not_authenticated')
        }

    def help_test_success(self, request_data):
        response = self.client.patch(self.url,
                                     request_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin.username,
                                     HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {'msg': 'Task update successful.',
                                         'id': self.task.id})
        self.assertEqual(self.task.id, response.data.get('id'))

    def test_title_success(self):
        request_data = {'id': self.task.id, 'data': {'title': 'New Title'}}
        self.help_test_success(request_data)
        self.assertEqual(Task.objects.get(id=self.task.id).title,
                         request_data.get('data').get('title'))

    def test_title_blank(self):
        request_data = {'id': self.task.id, 'data': {'title': ''}}
        response = self.client.patch(self.url,
                                     request_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin.username,
                                     HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'data.title': ErrorDetail(string='Task title cannot be empty.',
                                      code='blank')
        })
        self.assertEqual(Task.objects.get(id=self.task.id).title,
                         self.task.title)

    def test_order_success(self):
        request_data = {'id': self.task.id, 'data': {'order': 10}}
        self.help_test_success(request_data)
        self.assertEqual(Task.objects.get(id=self.task.id).order,
                         request_data.get('data').get('order'))

    def test_order_blank(self):
        request_data = {'id': self.task.id, 'data': {'order': ''}}
        response = self.client.patch(self.url,
                                     request_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin.username,
                                     HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'data.order': ErrorDetail(string='Task order cannot be empty.',
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
        request_data = {'id': self.task.id,
                        'data': {'column': another_column.id}}
        self.help_test_success(request_data)
        self.assertEqual(Task.objects.get(id=self.task.id).column.id,
                         request_data.get('data').get('column'))

    def test_column_blank(self):
        request_data = {'id': self.task.id, 'data': {'column': ''}}
        response = self.client.patch(self.url,
                                     request_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin.username,
                                     HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'data.column': ErrorDetail(string='Task column cannot be empty.',
                                       code='blank')
        })
        self.assertEqual(Task.objects.get(id=self.task.id).column,
                         self.task.column)

    def test_column_invalid(self):
        request_data = {'id': self.task.id, 'data': {'column': '123123'}}
        response = self.client.patch(self.url,
                                     request_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin.username,
                                     HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'data.column': ErrorDetail(string='Invalid column id.',
                                       code='invalid')
        })
        self.assertEqual(Task.objects.get(id=self.task.id).column,
                         self.task.column)

    def test_auth_token_empty(self):
        initial_count = Task.objects.count()
        request_data = {'id': self.task.id, 'data': {'title': 'New Title'}}
        response = self.client.post(self.url,
                                    request_data,
                                    format='json',
                                    HTTP_AUTH_USER=self.admin.username,
                                    HTTP_AUTH_TOKEN='')
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, self.forbidden_response)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_token_invalid(self):
        initial_count = Task.objects.count()
        request_data = {'id': self.task.id, 'data': {'title': 'New Title'}}
        response = self.client.post(self.url,
                                    request_data,
                                    format='json',
                                    HTTP_AUTH_USER=self.admin.username,
                                    HTTP_AUTH_TOKEN='ASDKFJ!FJ_012rjpiwajfosi')
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, self.forbidden_response)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_user_blank(self):
        initial_count = Task.objects.count()
        request_data = {'id': self.task.id, 'data': {'title': 'New Title'}}
        response = self.client.post(self.url,
                                    request_data,
                                    format='json',
                                    HTTP_AUTH_USER='',
                                    HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, self.forbidden_response)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_user_invalid(self):
        initial_count = Task.objects.count()
        request_data = {'id': self.task.id, 'data': {'title': 'New Title'}}
        response = self.client.post(self.url,
                                    request_data,
                                    format='json',
                                    HTTP_AUTH_USER='invalidio',
                                    HTTP_AUTH_TOKEN=self.admin_token)
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, self.forbidden_response)
        self.assertEqual(Board.objects.count(), initial_count)
