from rest_framework.test import APITestCase
from ..models import Board, Column, Team


class CreateColumnTests(APITestCase):
    def setUp(self):
        self.url = '/columns/'
        team = Team.objects.create()
        self.board = Board.objects.create(team=team)

    def test_success(self):
        initial_count = Column.objects.filter(board=self.board)
        response = self.client.post(self.url, {'board_id': self.board.id})
        self.assertEqual(response.status_code, 201)
        self.assertEqual(Column.objects.filter(board=self.board),
                         initial_count + 1)
