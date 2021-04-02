from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Board, Team, User, Column


class CreateBoardTests(APITestCase):
    def setUp(self):
        self.url = '/boards/'
        self.team = Team.objects.create()
        self.admin = User.objects.create(username='teamadmin',
                                         password='adminpassword',
                                         is_admin=True,
                                         team=self.team)
        self.member = User.objects.create(username='teammember',
                                          password='memberpassword',
                                          is_admin=False,
                                          team=self.team)

    def test_success(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url, {'username': self.admin.username,
                                               'team_id': self.team.id})
        self.assertEqual(response.status_code, 201)
        board_id = response.data.get('board_id')
        self.assertTrue(board_id)
        self.assertEqual(response.data.get('team_id'), self.team.id)
        self.assertEqual(Board.objects.count(), initial_count + 1)
        columns = Column.objects.filter(board=board_id)
        self.assertEqual(len(columns), 4)

    def test_username_blank(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url, {'username': '',
                                               'team_id': self.team.id})
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': ErrorDetail(string='Username cannot be empty.',
                                    code='blank')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_username_invalid(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url, {'username': 'invalidio',
                                               'team_id': self.team.id})
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': ErrorDetail(string='Invalid username.', code='invalid')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_user_not_admin(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url, {
            'username': self.member.username,
            'team_id': self.team.id
        })
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': ErrorDetail(
                string='Only the team admin can create a board.',
                code='not_authorized'
            )
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_team_id_blank(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url, {'username': self.admin.username,
                                               'team_id': ''})
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team ID cannot be empty.',
                                   code='blank')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_team_not_found(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url, {'username': self.admin.username,
                                               'team_id': '123'})
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team not found.', code='not_found')
        })
        self.assertEqual(Board.objects.count(), initial_count)


