from rest_framework.test import APITestCase
from ..models import Subtask, Task, Column, Board, Team


class UpdateSubtask(APITestCase):
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

    def test_update_title_success(self):
        request = {
            'id': self.subtask.id,
            'title': 'New Task Title'
        }
        response = self.client.patch(self.url, request)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': 'Subtask update successful.',
            'id': self.subtask.id
        })
        subtask = Subtask.objects.get(id=self.subtask.id)
        self.assertEqual(subtask.title, request.get('title'))
