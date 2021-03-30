from rest_framework.exceptions import ErrorDetail
from rest_framework.test import APITestCase
from main.models import User, Team
from uuid import uuid4


# noinspection DuplicatedCode
class RegisterTestCase(APITestCase):
    url = '/register/'

    def test_success(self):
        initial_user_count = User.objects.count()
        initial_team_count = Team.objects.count()
        request_data = {'username': 'fooooooooo',
                        'password': 'barbarbar',
                        'password_confirmation': 'barbarbar'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.data['username'], request_data['username'])
        self.assertEqual(response.data['password'], request_data['password'])
        self.assertTrue(response.data['is_admin'])
        self.assertEqual(User.objects.count(), initial_user_count + 1)
        self.assertEqual(Team.objects.count(), initial_team_count + 1)

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
        self.assertEqual(User.objects.count(), initial_count + 1)
        self.assertEqual(request_data['username'], response.data['username'])
        self.assertEqual(request_data['password'], response.data['password'])
        self.assertFalse(response.data['is_admin'])

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
        request_data = {'username': 'foo',
                        'password': 'bar',
                        'password_confirmation': 'bar',
                        'invite_code': invite_code}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, {'invite_code': "Team not found."})
        self.assertEqual(User.objects.count(), initial_user_count)
        self.assertEqual(Team.objects.count(), initial_team_count)

    def test_unmatched_passwords(self):
        initial_count = User.objects.count()
        response = self.client.post(self.url, {
            'username': 'foooo',
            'password_confirmation': 'barbarbar',
            'password': 'barbarbar'
        })
        self.assertEqual(response.status_code, 400)
        print(response.data)
        self.assertEqual(response.data, {
            'password_confirmation': "Confirmation does not match the "
                                     "password."
        })
        self.assertEqual(User.objects.count(), initial_count)

    def test_empty_username_field(self):
        initial_user_count = User.objects.count()
        initial_team_count = Team.objects.count()
        request_data = {'password': 'bar',
                        'password_confirmation': 'bar'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data['username'][0],
                         'This field may not be null.')
        self.assertEqual(User.objects.count(), initial_user_count)
        self.assertEqual(Team.objects.count(), initial_team_count)

    def test_empty_password_field(self):
        initial_user_count = User.objects.count()
        initial_team_count = Team.objects.count()
        request_data = {'username': 'foo',
                        'password_confirmation': 'bar'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'password': 'Password cannot be empty.'
        })
        self.assertEqual(User.objects.count(), initial_user_count)
        self.assertEqual(Team.objects.count(), initial_team_count)

    def test_empty_password_confirmation_field(self):
        initial_user_count = User.objects.count()
        initial_team_count = Team.objects.count()
        request_data = {'username': 'foo',
                        'password': 'bar'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'password_confirmation': 'Password confirmation cannot be empty.'
        })
        self.assertEqual(User.objects.count(), initial_user_count)
        self.assertEqual(Team.objects.count(), initial_team_count)
