from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
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
        initial_count = Task.objects.count()
        request = {'title': 'Some Task',
                   'description': 'Lorem ipsum dolor sit amet',
                   'column': self.column.id}
        response = self.client.post(self.url, request)
        self.assert_success(response.data, response.status_code, request)
        self.assertEqual(Task.objects.count(), initial_count + 1)

    def test_success_without_description(self):
        initial_count = Task.objects.count()
        request = {'title': 'Some Task',
                   'description': '',
                   'column': self.column.id}
        response = self.client.post(self.url, request)
        self.assert_success(response.data, response.status_code, request)
        self.assertEqual(Task.objects.count(), initial_count + 1)

    def test_success_with_subtasks(self):
        initial_count = Task.objects.count()
        request = {'title': 'Some Task',
                   'description': 'Lorem ipsum dolor sit amet',
                   'column': self.column.id,
                   'subtasks': [{'title': 'Do something'},
                                {'title': 'Do some other thing'}]}
        response = self.client.post(self.url, request, format='json')
        self.assert_success(response.data, response.status_code, request)
        subtasks = Subtask.objects.filter(task=response.data.get('task_id'))
        self.assertEqual(subtasks.count(), len(request.get('subtasks')))
        self.assertEqual(Task.objects.count(), initial_count + 1)

    def test_title_blank(self):
        initial_count = Task.objects.count()
        response = self.client.post(self.url, {
            'title': '',
            'description': 'Lorem ipsum dolor sit amet',
            'column': self.column.id
        })
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'title': [ErrorDetail(string='Title cannot be empty.',
                                  code='blank')]
        })
        self.assertEqual(Task.objects.count(), initial_count)

    def test_title_max_length(self):
        initial_count = Task.objects.count()
        response = self.client.post(self.url, {
            'title': 'foooooooooooooooooooooooooooooooooooooooooooooooooo',
            'description': 'Lorem ipsum dolor sit amet',
            'column': self.column.id
        })
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'title': [
                ErrorDetail(string='Title cannot be longer than 50 characters.',
                            code='max_length'),
            ]
        })
        self.assertEqual(Task.objects.count(), initial_count)

    def test_column_blank(self):
        initial_count = Task.objects.count()
        request = {'title': 'Some Task',
                   'description': 'Lorem ipsum dolor sit amet',
                   'column': ''}
        response = self.client.post(self.url, request)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'column': ErrorDetail(string='Column cannot be empty.',
                                  code='blank')
        })
        self.assertEqual(Task.objects.count(), initial_count)
