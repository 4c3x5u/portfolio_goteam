from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Board, Team, User


class ListBoardsTests(APITestCase):
    def setUp(self):
        team = Team.objects.create()
        for _ in range(0, 3):
            Board.objects.create(team_id=team.id)
        self.base_url = '/boards/'
        self.team_id = team.id

    def test_success(self):
        response = self.client.get(f'{self.base_url}?team_id={self.team_id}')
        print(f'§{response.data}')
        self.assertEqual(response.status_code, 200)
        self.assertTrue(response.data.get('boards'))

