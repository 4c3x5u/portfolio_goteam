from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from main.models import User, Team


# noinspection DuplicatedCode
class LoginTests(APITestCase):
    url = '/login/'

    def setUp(self):
        User.objects.create(username='foooo',
                            password='barbarbar',
                            team=Team.objects.create())

    def test_success(self):
        request_data = {'username': 'foooo', 'password': 'barbarbar'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': 'Login successful.',
            'username': 'foooo'
        })
        user = User.objects.get(username='foooo')
        self.assertTrue(user)
        self.assertEqual(user.password, request_data['password'])

    def test_username_blank(self):
        request_data = {'username': '', 'password': 'barbarbar'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': [
                ErrorDetail(string='Username cannot be empty.',
                            code='blank')
            ]
        })

    def test_password_blank(self):
        request_data = {'username': 'foooo', 'password': ''}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'password': [
                ErrorDetail(string='Password cannot be empty.',
                            code='blank')
            ]
        })

    def test_username_invalid(self):
        request_data = {'username': 'not_foooo', 'password': 'barbarbar'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': [
                ErrorDetail(string='Invalid username.',
                            code='invalid')
            ]
        })

    def test_password_invalid(self):
        request_data = {'username': 'foooo', 'password': 'not_barbarbar'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'password': [
                ErrorDetail(string='Invalid password.',
                            code='invalid')
            ]
        })
