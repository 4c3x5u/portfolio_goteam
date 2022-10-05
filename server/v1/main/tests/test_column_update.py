from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Column, Board, Team, Task
from main.helpers.user_helper import UserHelper
from main.helpers.auth_helper import AuthHelper


class UpdateColumnTests(APITestCase):
    endpoint = '/columns/?id='

    def setUp(self):
        team = Team.objects.create()
        board = Board.objects.create(name='My Board', team=team)
        self.column = Column.objects.create(order=0, board=board)
        self.another_column = Column.objects.create(order=1, board=board)
        self.tasks = [Task.objects.create(
            title=str(i), order=i, column=self.column
        ) for i in range(0, 5)]

        user_helper = UserHelper(team)
        self.member = user_helper.create_user()
        self.admin = user_helper.create_user(is_admin=True)
        self.assigned_member = user_helper.create_user()

        wrong_user_helper = UserHelper(Team.objects.create())
        self.wrong_admin = wrong_user_helper.create_user(is_admin=True)

        self.task_data = [{
            'id': task.id,
            'title': task.title,
            'order': 5 - task.order,
            'user': self.assigned_member['username']
        } for task in self.tasks]

    def help_test_success(self, user):
        response = self.client.patch(f'{self.endpoint}{self.column.id}',
                                     self.task_data,
                                     format='json',
                                     HTTP_AUTH_USER=user['username'],
                                     HTTP_AUTH_TOKEN=user['token'])
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': 'Column and all its tasks updated successfully.',
            'id': self.column.id,
        })
        new_tasks = Task.objects.filter(column_id=self.column.id)
        for i in range(0, 5):
            task = new_tasks.get(title=str(i))
            self.assertEqual(task.order, 5 - int(task.title))


    def test_admin_success(self):
        self.help_test_success(self.admin)

    def test_assigned_member_success(self):
        self.help_test_success(self.assigned_member)

    def test_column_id_empty(self):
        response = self.client.patch(self.endpoint,
                                     self.task_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'column': [ErrorDetail(string='Column ID cannot be null.',
                                   code='null')]
        })
        new_tasks = Task.objects.filter(column_id=self.column.id)
        for i in range(0, 5):
            task = new_tasks.get(title=str(i))
            self.assertEqual(task.order, int(task.title))

    def test_task_id_empty(self):
        request_data = list(map(
            lambda t: {'title': t['title'], 'order': t['order']},
            self.task_data
        ))
        response = self.client.patch(f'{self.endpoint}{self.column.id}',
                                     request_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'task.id': ErrorDetail(string='Task ID cannot be empty.',
                                   code='blank')
        })
        new_tasks = Task.objects.filter(column_id=self.column.id)
        for i in range(0, 5):
            task = new_tasks.get(title=str(i))
            self.assertEqual(task.order, int(task.title))

    def test_auth_token_empty(self):
        response = self.client.patch(f'{self.endpoint}{self.column.id}',
                                     self.task_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN='')
        self.assertEqual(response.status_code,
                         AuthHelper.AUTHENTICATION_ERROR.status_code)
        self.assertEqual(response.data,
                         AuthHelper.AUTHENTICATION_ERROR.detail)

    def test_auth_token_invalid(self):
        response = self.client.patch(f'{self.endpoint}{self.column.id}',
                                     self.task_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN='ASDKFJ!FJ_012rjpiwajfos')
        self.assertEqual(response.status_code,
                         AuthHelper.AUTHENTICATION_ERROR.status_code)
        self.assertEqual(response.data,
                         AuthHelper.AUTHENTICATION_ERROR.detail)

    def test_auth_user_blank(self):
        response = self.client.patch(
            f'{self.endpoint}{self.column.id}',
            self.task_data,
            format='json',
            HTTP_AUTH_USER='',
            HTTP_AUTH_TOKEN=self.admin['token']
        )
        self.assertEqual(response.status_code,
                         AuthHelper.AUTHENTICATION_ERROR.status_code)
        self.assertEqual(response.data,
                         AuthHelper.AUTHENTICATION_ERROR.detail)

    def test_auth_user_invalid(self):
        response = self.client.patch(f'{self.endpoint}{self.column.id}',
                                     self.task_data,
                                     format='json',
                                     HTTP_AUTH_USER='invalidio',
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code,
                         AuthHelper.AUTHENTICATION_ERROR.status_code)
        self.assertEqual(response.data,
                         AuthHelper.AUTHENTICATION_ERROR.detail)

    def test_wrong_team(self):
        response = self.client.patch(
            f'{self.endpoint}{self.column.id}',
            self.task_data,
            format='json',
            HTTP_AUTH_USER=self.wrong_admin['username'],
            HTTP_AUTH_TOKEN=self.wrong_admin['token']
        )
        self.assertEqual(response.status_code,
                         AuthHelper.AUTHORIZATION_ERROR.status_code)
        self.assertEqual(response.data,
                         AuthHelper.AUTHORIZATION_ERROR.detail)

    def test_not_authorized(self):
        response = self.client.patch(
            f'{self.endpoint}{self.another_column.id}',
            self.task_data,
            format='json',
            HTTP_AUTH_USER=self.member['username'],
            HTTP_AUTH_TOKEN=self.member['token']
        )
        self.assertEqual(response.status_code,
                         AuthHelper.AUTHORIZATION_ERROR.status_code)
        self.assertEqual(response.data,
                         AuthHelper.AUTHORIZATION_ERROR.detail)
