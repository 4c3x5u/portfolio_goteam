from rest_framework.test import APITestCase
from ..models import Task, Column, Board, Team


class GetTasksTests(APITestCase):
    endpoint = 'tasks/column_id='

    def setUp(self):
        team = Team.objects.create()
        board = Board.objects.create(team=team)
        self.column = Column.objects.create(order=0, board=board)
        self.tasks = list(map(
            lambda task: {
                'id': task.id,
                'title': task.title,
                'description': task.description,
                'order': task.order
            }, [
                Task.objects.create(
                    title=f'Task #{i}',
                    order=i,
                    column=self.column
                ) for i in range(0, 10)
            ]
        ))

    def test_success(self):
        response = self.client.get(f'{self.endpoint}{self.column.id}')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data.get('tasks'), self.tasks)
