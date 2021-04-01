from rest_framework.test import APITestCase
from ..models import Team, Board, Column


class CreateTaskTests(APITestCase):
    def setUp(self):
        self.url = '/tasks/'
        team = Team.objects.create()
        board = Board.objects.create(team=team)
        self.column = Column.objects.create(board=board)

    def test_success(self):
        request = {
            # TODO: Set order during the creation based on the pre-existing
            #       tasks inside the column
            'title': 'Some Task',
            'description': 'Lorem ipsum dolor sit amet',
            'column': self.column.id,
            'subtasks': ['Some Subtask', 'Some Other Subtask']
        }
        response = self.client.post(self.url, request)
        self.assertEqual(response.status_code, 201)
