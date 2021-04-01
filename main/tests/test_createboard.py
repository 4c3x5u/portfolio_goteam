from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Board, Team, User


class CreateBoardTests(APITestCase):
    def setUp(self):
        self.url = '/boards/'

    def test_success(self):
        initial_count = Board.objects.count()
        team = Team.objects.create()
        user = User.objects.create(username='foooo',
                                   password='barbarbar',
                                   is_admin=True,
                                   team=team)
        response = self.client.post(self.url, {'username': user.username,
                                               'team_id': team.id})
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.data.get('team_id'), team.id)
        self.assertTrue(response.data.get('board_id'))
        self.assertEqual(Board.objects.count(), initial_count + 1)

    def test_username_blank(self):
        initial_count = Board.objects.count()
        team = Team.objects.create()
        response = self.client.post(self.url, {'username': '',
                                               'team_id': team.id})
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': ErrorDetail(string='Username cannot be empty.',
                                    code='blank')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_username_invalid(self):
        initial_count = Board.objects.count()
        team = Team.objects.create()
        response = self.client.post(self.url, {'username': 'some_username',
                                               'team_id': team.id})
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': ErrorDetail(string='Invalid username.', code='invalid')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_user_not_admin(self):
        initial_count = Board.objects.count()
        team = Team.objects.create()
        user = User.objects.create(username='foooo',
                                   password='barbarbar',
                                   is_admin=False,
                                   team=team)
        response = self.client.post(self.url, {
            'username': user.username,
            'team_id': team.id
        })
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': ErrorDetail(
                string='Only the team admin can create a board.',
                code='not_authorized'
            )
        })
        self.assertEqual(Board.objects.count(), initial_count)
