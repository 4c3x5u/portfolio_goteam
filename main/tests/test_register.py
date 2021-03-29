from rest_framework.test import APITestCase
from main.models import User, Team
from uuid import uuid4


# noinspection DuplicatedCode
class RegisterTestCase(APITestCase):
    def test_success(self):
        initial_count = User.objects.count()
        request_data = {'username': 'foo',
                        'password': 'bar',
                        'password_confirmation': 'bar'}
        response = self.client.post('/user/', request_data)
        response.status_code != 201 and print(response.data)
        self.assertEqual(User.objects.count(), initial_count + 1)
        for attr, expected_value in request_data.items():
            if attr != 'password_confirmation' and attr != 'invite_code':
                self.assertEqual(response.data[attr], expected_value)
        self.assertTrue(Team.objects.get(id=response.data['team']))
        self.assertTrue(response.data['is_admin'])

    def test_success_with_invite_code(self):
        initial_count = User.objects.count()
        team = Team.objects.create()
        ic = team.invite_code
        request_data = {'username': 'foo',
                        'password': 'bar',
                        'password_confirmation': 'bar',
                        'invite_code': ic}
        response = self.client.post('/user/', request_data)
        response.status_code != 201 and print(response.data)
        self.assertEqual(User.objects.count(), initial_count + 1)
        for attr, expected_value in request_data.items():
            if attr != 'password_confirmation' and attr != 'invite_code':
                self.assertEqual(response.data[attr], expected_value)
        self.assertEqual(response.data['team'], team.id)
        self.assertFalse(response.data['is_admin'])

    def test_team_not_found(self):
        initial_user_count = User.objects.count()
        initial_team_count = Team.objects.count()
        invite_code = uuid4()
        request_data = {'username': 'foo',
                        'password': 'bar',
                        'password_confirmation': 'bar',
                        'invite_code': invite_code}
        response = self.client.post('/user/', request_data)
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, {'invite_code': "team not found"})
        self.assertEqual(User.objects.count(), initial_user_count)
        self.assertEqual(Team.objects.count(), initial_team_count)

    def test_unmatched_passwords(self):
        request_data = {'username': 'foo',
                        'password': 'bar',
                        'password_confirmation': 'not_bar'}
        initial_count = User.objects.count()
        response = self.client.post('/user/', request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'password_confirmation': "confirmation doesn't match the password"
        })
        self.assertEqual(User.objects.count(), initial_count)
