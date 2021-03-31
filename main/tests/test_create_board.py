from rest_framework.test import APITestCase
from ..models import Board, Team, User


class CreateBoardTests(APITestCase):
    def setUp(self):
        self.url = '/board/'
        self.team = Team.objects.create()
        self.user = User.objects.create(username='foooo',
                                        password='barbarbar',
                                        is_admin=True,
                                        team=self.team)
        self.initial_board_count = Board.objects.count()

    def test_success(self):
        request_data = {'username': self.user.username}
        response = self.client.post(self.url, request_data)
        board = Board.objects.get(team=self.team)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.data, {
            'msg': 'Board created successfuly',
            'team_id': self.team.id,
            'board_id': board.id
        })
        self.assertEqual(Board.objects.count(), self.initial_board_count + 1)
        self.assertEqual(board.team, self.team)
