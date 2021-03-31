from rest_framework.exceptions import ErrorDetail
from rest_framework.test import APITestCase
from main.models import User, Team, Board
from uuid import uuid4


# noinspection DuplicatedCode
class RegisterTests(APITestCase):
    url = '/register/'

    def test_success(self):
        initial_user_count = User.objects.count()
        initial_team_count = Team.objects.count()
        request_data = {'username': 'fooooooooo',
                        'password': 'barbarbar',
                        'password_confirmation': 'barbarbar'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.data, {
            'msg': 'Login successful.',
            'username': request_data['username'],
            'team_id': response.data['team_id']
        })
        self.assertEqual(User.objects.count(), initial_user_count + 1)
        self.assertEqual(Team.objects.count(), initial_team_count + 1)
        self.assertEqual(
            Board.objects.filter(team=response.data['team_id']).count(),
            1
        )

    def test_success_with_invite_code(self):
        initial_count = User.objects.count()
        team = Team.objects.create()
        ic = team.invite_code
        request_data = {'username': 'foooo',
                        'password': 'barbarbar',
                        'password_confirmation': 'barbarbar',
                        'invite_code': ic}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.data, {
            'msg': 'Login successful.',
            'username': request_data['username'],
            'team_id': team.id
        })
        self.assertEqual(User.objects.count(), initial_count + 1)

    def test_invalid_invite_code(self):
        initial_user_count = User.objects.count()
        initial_team_count = Team.objects.count()
        request_data = {'username': 'foooo',
                        'password': 'barbarbar',
                        'password_confirmation': 'barbarbar',
                        'invite_code': 'invalid uuid'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(
            response.data,
            {
                'invite_code': [
                    ErrorDetail(string='Invalid invite code.', code='invalid')
                ]
            }
        )
        self.assertEqual(User.objects.count(), initial_user_count)
        self.assertEqual(Team.objects.count(), initial_team_count)

    def test_team_not_found(self):
        initial_user_count = User.objects.count()
        initial_team_count = Team.objects.count()
        invite_code = uuid4()
        request_data = {'username': 'foooo',
                        'password': 'barbarbar',
                        'password_confirmation': 'barbarbar',
                        'invite_code': invite_code}
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
        print(response.data)
        self.assertEqual(response.data, {
            'password_confirmation': [
                ErrorDetail(string='Confirmation does not match the password.',
                            code='invalid')
            ]
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
            'password_confirmation': [
                ErrorDetail(
                    string='Password confirmation cannot be empty.',
                    code='blank'
                )
            ]
        })
        self.assertEqual(User.objects.count(), initial_user_count)
        self.assertEqual(Team.objects.count(), initial_team_count)
