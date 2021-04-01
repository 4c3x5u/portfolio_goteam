from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Board, Team, User


class ListBoardsTests(APITestCase):
    def setUp(self):
        self.team = Team.objects.create()
        self.boards = []
        for _ in range(0, 3):
            self.boards.append(Board.objects.create(team_id=self.team.id))
        self.base_url = '/boards/'
        self.team_id = self.team.id

    def test_success(self):
        response = self.client.get(f'{self.base_url}?team_id={self.team_id}')
        self.assertEqual(response.status_code, 200)
        boards = response.data.get('boards')
        self.assertTrue(boards)
        self.assertTrue(boards.count, 3)
        for board in boards:
            self.assertEqual(board.get('team_id'), self.team.id)
