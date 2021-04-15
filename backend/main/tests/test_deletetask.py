from rest_framework.test import APITestCase
from ..models import Team, Board, Column, Task


class DeleteTaskTests(APITestCase):
    endpoint = '/tasks/?id='

    def setUp(self):
        team = Team.objects.create()
        board = Board.objects.create(team=team)
        column = Column.objects.create(order=0, board=board)
        self.task = Task.objects.create(
            title='Do Something!',
            order=0,
            column=column
        )

    def test_success(self):
        initial_count = Task.objects.count()
        response = self.client.delete(f'{self.endpoint}{self.task.id}')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': 'Task deleted successfully.',
            'id': self.task.id,
        })
        self.assertEqual(Task.objects.count(), initial_count - 1)
