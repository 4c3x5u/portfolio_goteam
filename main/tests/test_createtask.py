from rest_framework.test import APITestCase
from ..models import Team, Board, Column, Task, Subtask


class CreateTaskTests(APITestCase):
    def setUp(self):
        self.url = '/tasks/'
        team = Team.objects.create()
        board = Board.objects.create(team=team)
        self.column = Column.objects.create(board=board, order=0)

    def test_success(self):
        request = {'title': 'Some Task',
                   'description': 'Lorem ipsum dolor sit amet',
                   'column': self.column.id,
                   'subtasks': [{'title': 'Do something'},
                                {'title': 'Do some other thing'}]}
        response = self.client.post(self.url, request, format='json')
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.data.get('msg'), 'Task creation successful.')
        task_id = response.data.get('task_id')
        self.assertTrue(task_id)
        task = Task.objects.get(id=task_id)
        self.assertEqual(task.title, request.get('title'))
        self.assertEqual(task.description, request.get('description'))
        self.assertEqual(task.column.id, request.get('column'))
        subtasks = Subtask.objects.filter(task=task_id)
        self.assertEqual(subtasks.count(), len(request.get('subtasks')))
