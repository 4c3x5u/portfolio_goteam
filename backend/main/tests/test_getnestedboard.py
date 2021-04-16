from rest_framework.test import APITestCase
from ..models import Board, Team, Column, Task, Subtask
from ..util import new_member


class GetNestedBoardTests(APITestCase):
    endpoint = '/boards/?board_id='

    def setUp(self):
        self.team = Team.objects.create()
        self.member = new_member(self.team)
        self.boards = [
            Board.objects.create(team_id=self.team.id) for _ in range(0, 3)
        ]
        self.columns = [
            Column.objects.create(
                board_id=self.boards[0].id, order=i
            ) for i in range(0, 4)
        ]
        self.tasks = [
            Task.objects.create(
                title=f'Task #{i}',
                order=i,
                column=self.columns[0]
            ) for i in range(0, 5)
        ]
        self.subtasks = [
            Subtask.objects.create(
                title=f'Subtask #{i}',
                order=i,
                task=self.tasks[0],
                done=(i % 2 == 0)
            ) for i in range(0, 2)
        ]

    def test_nested_success(self):
        response = self.client.get(f'{self.endpoint}{self.boards[0].id}',
                                   HTTP_AUTH_USER=self.member['username'],
                                   HTTP_AUTH_TOKEN=self.member['token'])
        self.assertEqual(response.status_code, 200)
        self.assertTrue(response.data.get('id'), self.boards[0].id)

        columns = response.data.get('columns')
        self.assertEqual(len(columns), 4)
        for i in range(0, 4):
            self.assertEqual(columns[i].get('id'), self.columns[i].id)
            self.assertEqual(columns[i].get('order'), self.columns[i].order)

        tasks = columns[0].get('tasks')
        self.assertEqual(len(tasks), 5)
        for i in range(0, 5):
            self.assertEqual(tasks[i].get('id'), self.tasks[i].id)
            self.assertEqual(tasks[i].get('title'), self.tasks[i].title)
            self.assertEqual(tasks[i].get('description'),
                             self.tasks[i].description)
            self.assertEqual(tasks[i].get('order'), self.tasks[i].order)

        subtasks = tasks[0].get('subtasks')
        self.assertEqual(len(subtasks), 2)
        for i in range(0, 2):
            self.assertEqual(subtasks[i].get('id'), self.subtasks[i].id)
            self.assertEqual(subtasks[i].get('title'), self.subtasks[i].title)
            self.assertEqual(subtasks[i].get('order'), self.subtasks[i].order)
            self.assertEqual(subtasks[i].get('done'), self.subtasks[i].done)

