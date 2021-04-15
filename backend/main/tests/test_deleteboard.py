from rest_framework.test import APITestCase
from ..models import Team, Board
from ..util import new_admin

class DeleteBoardTests(APITestCase):
    endpoint = '/boards/?id='

    def setUp(self):
        team = Team.objects.create()
        self.board = Board.objects.create(team=team)
        self.admin = new_admin(team)

    def test_success(self):
        initial_count = Board.objects.count()
        response = self.client.delete(f'{self.endpoint}{self.board.id}',
                                      HTTP_AUTH_USER=self.admin.username,
                                      HTTP_AUTH_TOKEN=self.admin.token)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': 'Board deleted successfully',
            'id': str(self.board.id),
        })
        self.assertEqual(Board.objects.count(), initial_count - 1)
