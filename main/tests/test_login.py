from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from main.models import User, Team


# noinspection DuplicatedCode
class LoginTestCase(APITestCase):
    url = '/login/'

    def setUp(self):
        user = User.objects.create(username='foooo',
                                   password='barbarbar',
                                   team=Team.objects.create())
        self.assertTrue(user)

    def test_success(self):
        request_data = {'username': 'foooo', 'password': 'barbarbar'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {'foooo': 'Login successful.'})
        user = User.objects.get(username='foooo')
        self.assertTrue(user)
        self.assertEqual(user.password, request_data['password'])

    def test_username_blank(self):
        request_data = {'username': '', 'password': 'barbarbar'}
        response = self.client.post(self.url, request_data)
        print(f'§{response.data}')
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': [
                ErrorDetail(string='Username cannot be empty.',
                            code='blank')
            ]
        })
        user = User.objects.get(username='foooo')
        self.assertTrue(user)
        self.assertEqual(user.password, request_data['password'])

