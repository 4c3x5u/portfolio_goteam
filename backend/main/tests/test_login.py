from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from main.models import User, Team


class LoginTests(APITestCase):
    def setUp(self):
        self.url = '/login/'
        self.pw_raw = 'barbarbar'
        self.user = User.objects.create(
            username='foooo',
            password=b'$2b$12$ZC.GGCmSPi8syzmJBQ6LoeUSeD2wkdSBZkPh18nZU81Lv6u7'
                     b'CuZMe',
            team=Team.objects.create()
        )

    def test_success(self):
        request_data = {'username': self.user.username,
                        'password': self.pw_raw}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data.get('msg'), 'Login successful.')
        self.assertEqual(response.data.get('username'), self.user.username)
        self.assertTrue(response.data.get('token'))

    def test_username_blank(self):
        request_data = {'username': '', 'password': self.pw_raw}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': [ErrorDetail(string='Username cannot be empty.',
                                     code='blank')]
        })

    def test_username_max_length(self):
        request_data = {'username': 'fooooooooooooooooooooooooooooooooooo',
                        'password': self.pw_raw}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': [
                ErrorDetail(
                    string='Username cannot be longer than 35 characters.',
                    code='max_length'
                )
            ]
        })

    # noinspection DuplicatedCode
    def test_password_max_length(self):
        password = '''
            barbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarb
            arbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarba
            rbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarbar
            barbarbarbarbarbarbarbarbarbarbarbarbarbarbarbarba
        '''
        request_data = {'username': self.user.username, 'password': password}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'password': [
                ErrorDetail(
                    string='Password cannot be longer than 255 characters.',
                    code='max_length'
                )
            ]
        })

    def test_password_blank(self):
        request_data = {'username': self.user.username, 'password': ''}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'password': [ErrorDetail(string='Password cannot be empty.',
                                     code='blank')]
        })

    def test_username_invalid(self):
        request_data = {'username': 'invalidusername',
                        'password': self.pw_raw}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'username': ErrorDetail(string='Invalid username.', code='invalid')
        })

    def test_password_invalid(self):
        request_data = {'username': self.user.username,
                        'password': 'invalidpassword'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'password': ErrorDetail(string='Invalid password.', code='invalid')
        })
