from rest_framework.test import APITestCase
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

    def test_update_title_success(self):
        request = {'id': self.subtask.id, 'data': {'title': 'New Task Title'}}
        subtask = self.help_test_success(request)
        self.assertEqual(subtask.title, request.get('data').get('title'))

    def test_update_done_success(self):
        request = {'id': self.subtask.id, 'data': {'done': True}}
        subtask = self.help_test_success(request)
        self.assertEqual(subtask.done, request.get('data').get('done'))

