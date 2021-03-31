from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Board, Team, User


class CreateBoardTests(APITestCase):
    def setUp(self):
        self.url = '/board/'
        self.team = Team.objects.create()
        self.user = User.objects.create(username='foooo',
                                        password='barbarbar',
                                        is_admin=True,
                                        team=self.team)

    def test_success(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url, {'username': self.user.username})
        board = Board.objects.get(team=self.team)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.data, {
            'msg': 'Board created successfuly',
            'team_id': self.team.id,
            'board_id': board.id
        })
        self.assertEqual(Board.objects.count(), initial_count + 1)
        self.assertEqual(board.team, self.team)

    def test_username_invalid(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url, {'username': 'some_username'})
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': [
                ErrorDetail(string='Invalid username.', code='invalid')
            ]
        })
        self.assertEqual(Board.objects.count(), initial_count)
