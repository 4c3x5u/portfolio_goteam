from rest_framework.exceptions import ErrorDetail
from rest_framework.test import APITestCase
from main.models import User, Team, Board
from uuid import uuid4


# noinspection DuplicatedCode
class RegisterTests(APITestCase):
    def setUp(self):
        self.url = '/register/'
        team = Team.objects.create()
        Board.objects.create(team=team)
        self.valid_invite_code = team.invite_code

    def test_success(self):
        initial_user_count = User.objects.count()
        initial_team_count = Team.objects.count()
        initial_board_count = Board.objects.count()
        request_data = {'username': 'fooooooooo',
                        'password': 'barbarbar',
                        'password_confirmation': 'barbarbar'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.data, {
            'msg': 'Registration successful.',
            'username': request_data['username'],
        })
        self.assertEqual(User.objects.count(), initial_user_count + 1)
        self.assertEqual(Team.objects.count(), initial_team_count + 1)
        self.assertEqual(Board.objects.count(), initial_board_count + 1)
        user = User.objects.get(username=response.data['username'])
        self.assertTrue(user)
        team = Team.objects.get(user=user)
        self.assertTrue(team)
        board = Board.objects.get(team=team)
        self.assertTrue(board)

    def test_success_with_invite_code(self):
        initial_count = User.objects.count()
        request_data = {'username': 'foooo',
                        'password': 'barbarbar',
                        'password_confirmation': 'barbarbar',
                        'invite_code': self.valid_invite_code}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.data, {
            'msg': 'Registration successful.',
            'username': request_data['username'],
        })
        self.assertEqual(User.objects.count(), initial_count + 1)
        user = User.objects.get(username=response.data['username'])
        self.assertTrue(user)
        team = Team.objects.get(user=user)
        self.assertTrue(team)
        board = Board.objects.get(team=team)
        self.assertTrue(board)

    def test_invalid_invite_code(self):
        initial_user_count = User.objects.count()
        initial_team_count = Team.objects.count()
        request_data = {'username': 'foooo',
                        'password': 'barbarbar',
                        'password_confirmation': 'barbarbar',
                        'invite_code': 'invalid uuid'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'invite_code': [
                ErrorDetail(string='Invalid invite code.', code='invalid')
            ]
        })
        self.assertEqual(User.objects.count(), initial_user_count)
        self.assertEqual(Team.objects.count(), initial_team_count)

    def test_team_not_found(self):
        initial_user_count = User.objects.count()
        initial_team_count = Team.objects.count()
        request_data = {'username': 'foooo',
                        'password': 'barbarbar',
                        'password_confirmation': 'barbarbar',
                        'invite_code': uuid4()}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'invite_code': [
                ErrorDetail(string='Team not found.', code='invalid')
            ]
        })
        self.assertEqual(User.objects.count(), initial_user_count)
        self.assertEqual(Team.objects.count(), initial_team_count)

    def test_unmatched_passwords(self):
        initial_count = User.objects.count()
        response = self.client.post(self.url, {
            'username': 'foooo',
            'password_confirmation': 'barbarbar',
            'password': 'not_barbarbar'
        })
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'password_confirmation': ErrorDetail(
                string='Confirmation does not match the password.',
                code='no_match'
            )
        })

        self.assertEqual(User.objects.count(), initial_count)

    def test_username_blank(self):
        initial_user_count = User.objects.count()
        initial_team_count = Team.objects.count()
        request_data = {'username': '',
                        'password': 'barbarbar',
                        'password_confirmation': 'barbarbar'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': [
                ErrorDetail(
                    string='Username cannot be empty.',
                    code='blank'
                )
            ]
        })
        self.assertEqual(User.objects.count(), initial_user_count)
        self.assertEqual(Team.objects.count(), initial_team_count)

    def test_password_blank(self):
        initial_user_count = User.objects.count()
        initial_team_count = Team.objects.count()
        request_data = {'username': 'foooo',
                        'password': '',
                        'password_confirmation': 'barbarbar'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'password': [
                ErrorDetail(
                    string='Password cannot be empty.',
                    code='blank',
                )
            ]
        })
        self.assertEqual(User.objects.count(), initial_user_count)
        self.assertEqual(Team.objects.count(), initial_team_count)

    def test_password_confirmation_blank(self):
        initial_user_count = User.objects.count()
        initial_team_count = Team.objects.count()
        request_data = {'username': 'foooo',
                        'password': 'barbarbar',
                        'password_confirmation': ''}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'password_confirmation': ErrorDetail(
                string='Password confirmation cannot be empty.',
                code='blank'
            )
        })
        self.assertEqual(User.objects.count(), initial_user_count)
        self.assertEqual(Team.objects.count(), initial_team_count)
