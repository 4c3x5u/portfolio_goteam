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
            'token': self.validToken
        })
