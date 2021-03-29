from rest_framework.test import APITestCase
from .models import User, Team


# noinspection DuplicatedCode
class RegisterTestCase(APITestCase):
    def test_valid_request(self):
        request_data = {'username': 'foo',
                        'password': 'bar',
                        'password_confirmation': 'bar'}
        initial_count = User.objects.count()
        response = self.client.post('/user/', request_data)
        response.status_code != 201 and print(response.data)
        self.assertEqual(User.objects.count(), initial_count + 1)
        for attr, expected_value in request_data.items():
            if attr != 'password_confirmation' and attr != 'invite_code':
                self.assertEqual(response.data[attr], expected_value)
        self.assertTrue(Team.objects.get(id=response.data['team']))
        self.assertTrue(response.data['is_admin'])
