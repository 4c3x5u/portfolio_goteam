from rest_framework.test import APITestCase
from main.models import User, Team


class VerifyTokenTests(APITestCase):
    def setUp(self):
        self.url = '/verify-token/'
        self.user = User.objects.create(
            username='foooo',
            password=b'$2b$12$S3tnLg/FbWsBAXaLtjk5nu6OwHWKs2spyMiI9W/.Kl/2Uh/j'
                     b'afyFC',
            team=Team.objects.create()
        )
        self.validToken = '$2b$12$78SqmRy1azwzJPnIgsWqoOBliZnuNr81oZZkcoDS8b' \
                          'B7TwvXYWHGq'

    def test_success(self):
        request_data = {'username': self.user.username,
                        'token': self.validToken}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': 'Token verification success.',
            'username': self.user.username,
            'teamId': self.user.team_id,
            'isAdmin': self.user.is_admin,
        })

    def test_token_invalid(self):
        request_data = {'username': self.user.username,
                        'token': 'as/dlkfjAS:DFkjaSdlnflasdjnvkasdjfasd,fasbd'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {'msg': 'Token verification failure.'})

    def test_token_empty(self):
        request_data = {'username': self.user.username,
                        'token': ''}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {'msg': 'Token verification failure.'})

    def test_username_invalid(self):
        request_data = {'username': 'nonexistent',
                        'token': 'as/dlkfjAS:DFkjaSdlnflasdjnvkasdjfasd,fasbd'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {'msg': 'Token verification failure.'})

    def test_username_empty(self):
        request_data = {'username': '',
                        'token': 'as/dlkfjAS:DFkjaSdlnflasdjnvkasdjfasd,fasbd'}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {'msg': 'Token verification failure.'})
