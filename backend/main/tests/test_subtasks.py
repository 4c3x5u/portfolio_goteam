from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Subtask, Task, Column, Board, Team


class SubtaskTests(APITestCase):
    def setUp(self):
        self.url = '/subtasks/'
        self.subtask = Subtask.objects.create(
            title='Some Task Title',
            order=0,
            task=Task.objects.create(
                title="Some Subtask Title",
                order=0,
                column=Column.objects.create(
                    order=0,
                    board=Board.objects.create(
                        team=Team.objects.create()
                    )
                )
            )
        )

    def help_test_success(self, request):
        response = self.client.patch(self.url, request, format='json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {'msg': 'Subtask update successful.',
                                         'id': self.subtask.id})
        return Subtask.objects.get(id=self.subtask.id)

    def help_test_failure(self):
        subtask = Subtask.objects.get(id=self.subtask.id)
        self.assertEqual(subtask.title, self.subtask.title)
        self.assertEqual(subtask.done, self.subtask.done)

    def test_update_title_success(self):
        request = {'id': self.subtask.id, 'data': {'title': 'New Task Title'}}
        subtask = self.help_test_success(request)
        self.assertEqual(subtask.title, request.get('data').get('title'))

    def test_update_done_success(self):
        request = {'id': self.subtask.id, 'data': {'done': True}}
        subtask = self.help_test_success(request)
        self.assertEqual(subtask.done, request.get('data').get('done'))

    def test_id_blank(self):
        request = {'id': '', 'data': {'title': 'New task title.'}}
        response = self.client.patch(self.url, request, format='json')
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'id': ErrorDetail(string='Subtask ID cannot be empty.',
                              code='blank')
        })
        self.help_test_failure()

    def test_data_blank(self):
        request = {'id': self.subtask.id, 'data': ''}
        response = self.client.patch(self.url, request, format='json')
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'data': ErrorDetail(string='Data cannot be empty.', code='blank')
        })
        self.help_test_failure()

    def test_title_blank(self):
        request = {'id': self.subtask.id, 'data': {'title': ''}}
        response = self.client.patch(self.url, request, format='json')
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'data.title': ErrorDetail(string='Title cannot be empty.',
                                      code='blank')
        })
        self.help_test_failure()

    def test_done_blank(self):
        request = {'id': self.subtask.id, 'data': {'done': ''}}
        response = self.client.patch(self.url, request, format='json')
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'data.done': ErrorDetail(string='Done cannot be empty.',
                                     code='blank')
        })
        self.help_test_failure()

