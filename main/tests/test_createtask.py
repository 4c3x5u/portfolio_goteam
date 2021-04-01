from rest_framework.test import APITestCase
from ..models import Team, Board, Column, Task, Subtask


class CreateTaskTests(APITestCase):
    def setUp(self):
        self.url = '/tasks/'
        team = Team.objects.create()
        board = Board.objects.create(team=team)
        self.column = Column.objects.create(board=board, order=0)

    def assert_success(self, response_data, status_code, request):
        self.assertEqual(status_code, 201)
        self.assertEqual(response_data.get('msg'), 'Task creation successful.')
        task_id = response_data.get('task_id')
        self.assertTrue(task_id)
        task = Task.objects.get(id=task_id)
        self.assertEqual(task.title, request.get('title'))
        self.assertEqual(task.description, request.get('description'))
        self.assertEqual(task.column.id, request.get('column'))

    def test_success(self):
        request = {'title': 'Some Task',
                   'description': 'Lorem ipsum dolor sit amet',
                   'column': self.column.id}
        response = self.client.post(self.url, request)
        self.assert_success(response.data, response.status_code, request)

    def test_success_with_subtasks(self):
        request = {'title': 'Some Task',
                   'description': 'Lorem ipsum dolor sit amet',
                   'column': self.column.id,
                   'subtasks': [{'title': 'Do something'},
                                {'title': 'Do some other thing'}]}
        response = self.client.post(self.url, request, format='json')
        self.assert_success(response.data, response.status_code, request)
        subtasks = Subtask.objects.filter(task=response.data.get('task_id'))
        self.assertEqual(subtasks.count(), len(request.get('subtasks')))
