from rest_framework.test import APITestCase


# noinspection DuplicatedCode
class LoginTestCase(APITestCase):
    register_url = '/register/'
    login_url = '/login/'

    def test_success(self):
        # Create user before trying to login
        self.client.post(self.register_url, {
            'username': 'foooo',
            'password': 'barbarbar',
            'password_confirmation': 'barbarbar'
        })

        request_data = {'username': 'foooo',
                        'password': 'barbarbar'}
        response = self.client.post(self.login_url, request_data)
        self.assertEqual(response.status_code, 200)

