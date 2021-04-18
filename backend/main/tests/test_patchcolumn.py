from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Column, Board, Team, Task
from ..util import new_admin, new_member, not_authenticated_response


class UpdateColumns(APITestCase):
    endpoint = '/columns/?id='

    def setUp(self):
        team = Team.objects.create()
        board = Board.objects.create(name='My Board', team=team)
        self.column = Column.objects.create(order=0, board=board)
        self.tasks = [
            Task.objects.create(
                title=str(i),
                order=i,
                column=self.column
            ) for i in range(0, 5)
        ]
        self.admin = new_admin(team)
        self.member = new_member(team)
        self.task_data = list(map(
            lambda task: {
                'id': task.id,
                'title': task.title,
                'order': 5 - task.order
            },
            self.tasks
        ))

    def test_success(self):
        response = self.client.patch(f'{self.endpoint}{self.column.id}',
                                     self.task_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': 'Column and all its tasks updated successfully.',
            'id': self.column.id,
        })
        new_tasks = Task.objects.filter(column_id=self.column.id)
        for i in range(0, 5):
            task = new_tasks.get(title=str(i))
            self.assertEqual(task.order, 5 - int(task.title))

    def test_column_id_empty(self):
        response = self.client.patch(self.endpoint,
                                     self.task_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'id': ErrorDetail(string='Column ID cannot be empty.',
                              code='blank')
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
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_token_invalid(self):
        response = self.client.patch(f'{self.endpoint}{self.column.id}',
                                     self.task_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.admin['username'],
                                     HTTP_AUTH_TOKEN='ASDKFJ!FJ_012rjpiwajfos')
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_user_blank(self):
        response = self.client.patch(
            f'{self.endpoint}{self.column.id}',
            self.task_data,
            format='json',
            HTTP_AUTH_USER='',
            HTTP_AUTH_TOKEN=self.admin['token']
        )
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_auth_user_invalid(self):
        response = self.client.patch(f'{self.endpoint}{self.column.id}',
                                     self.task_data,
                                     format='json',
                                     HTTP_AUTH_USER='invalidio',
                                     HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response.data)

    def test_not_admin(self):
        response = self.client.patch(f'{self.endpoint}{self.column.id}',
                                     self.task_data,
                                     format='json',
                                     HTTP_AUTH_USER=self.member['username'],
                                     HTTP_AUTH_TOKEN=self.member['token'])
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, {
            'auth': ErrorDetail(string='You must be an admin to do this.',
                                code='not_authorized')
        })

