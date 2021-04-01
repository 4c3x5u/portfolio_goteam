from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Board, Team, User


class ListBoardsTests(APITestCase):
    def setUp(self):
        self.base_url = '/boards/?team_id='
        self.team = Team.objects.create()
        self.boards = [
            Board.objects.create(team_id=self.team.id) for _ in range(0, 3)
        ]
        self.team_id = str(self.team.id)

    def test_success(self):
        response = self.client.get(self.base_url + self.team_id)
        self.assertEqual(response.status_code, 200)
        boards = response.data.get('boards')
        self.assertTrue(boards)
        self.assertTrue(boards.count, 3)
        for board in boards:
            self.assertEqual(board.get('team'), self.team.id)

    def test_boards_not_found(self):
        initial_count = Board.objects.count()
        team = Team.objects.create()
        response = self.client.get(self.base_url + str(team.id))
        self.assertEqual(response.status_code, 201)
        self.assertEqual(len(response.data.get('boards')), 1)
        self.assertEqual(Board.objects.count(), initial_count + 1)

    def test_team_id_empty(self):
        response = self.client.get(self.base_url)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team ID cannot be empty.',
                                   code='null')
        })

    def test_invalid_team_id(self):
        response = self.client.get(self.base_url + '123')
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team not found.', code='not_found')
        })
