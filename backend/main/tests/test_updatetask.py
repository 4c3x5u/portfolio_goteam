from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Task, Column, Board, Team


class UpdateTaskTests(APITestCase):
    def setUp(self):
        self.url = '/tasks/'
        self.task = Task.objects.create(
            title="Task Title",
            order=0,
            column=Column.objects.create(
                order=0,
                board=Board.objects.create(
                    team=Team.objects.create()
                )
            )
        )

    def help_test_success(self, request):
        response = self.client.patch(self.url, request, format='json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': 'Task update successful.',
            'id': self.task.id
        })
        # TODO: Implement something similar to tests all over
        self.assertEqual(self.task.id, response.data.get('id'))

    def test_title_success(self):
        request = {'id': self.task.id, 'data': {'title': 'New Title'}}
        self.help_test_success(request)
        self.assertEqual(Task.objects.get(id=self.task.id).title,
                         request.get('data').get('title'))

    def test_title_blank(self):
        request = {'id': self.task.id, 'data': {'title': ''}}
        response = self.client.patch(self.url, request, format='json')
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'data.title': ErrorDetail(string='Task title cannot be empty.',
                                      code='blank')
        })
        self.assertEqual(Task.objects.get(id=self.task.id).title,
                         self.task.title)

    def test_order_success(self):
        request = {'id': self.task.id, 'data': {'order': 10}}
        self.help_test_success(request)
        self.assertEqual(Task.objects.get(id=self.task.id).order,
                         request.get('data').get('order'))

    def test_order_blank(self):
        request = {'id': self.task.id, 'data': {'order': ''}}
        response = self.client.patch(self.url, request, format='json')
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
        request = {'id': self.task.id, 'data': {'column': another_column.id}}
        self.help_test_success(request)
        self.assertEqual(Task.objects.get(id=self.task.id).column.id,
                         request.get('data').get('column'))

    def test_column_blank(self):
        request = {'id': self.task.id, 'data': {'column': ''}}
        response = self.client.patch(self.url, request, format='json')
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'data.column': ErrorDetail(string='Task column cannot be empty.',
                                       code='blank')
        })
        self.assertEqual(Task.objects.get(id=self.task.id).column,
                         self.task.column)

    def test_column_invalid(self):
        request = {'id': self.task.id, 'data': {'column': '123123'}}
        response = self.client.patch(self.url, request, format='json')
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'data.column': ErrorDetail(string='Invalid column id.',
                                       code='invalid')
        })
        self.assertEqual(Task.objects.get(id=self.task.id).column,
                         self.task.column)
