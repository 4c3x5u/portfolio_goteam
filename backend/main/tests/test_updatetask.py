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

    def test_title_success(self):
        request = {'id': self.task.id, 'data': {'title': 'New Title'}}
        response = self.client.patch(self.url, request, format='json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': 'Task update successful.',
            'id': self.task.id
        })
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
        response = self.client.patch(self.url, request, format='json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': 'Task update successful.',
            'id': self.task.id
        })
        self.assertEqual(Task.objects.get(id=self.task.id).order,
                         request.get('data').get('order'))

