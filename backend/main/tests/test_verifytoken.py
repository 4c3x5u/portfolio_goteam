from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from main.models import User, Team
import bcrypt


class LoginTests(APITestCase):
    def setUp(self):
        self.url = '/verify-token/'
        self.user = User.objects.create(
            username='foooo',
            password=b'$2b$12$ZC.GGCmSPi8syzmJBQ6LoeUSeD2wkdSBZkPh18nZU81Lv6u7'
                     b'CuZMe',
            team=Team.objects.create()
        )
        self.validToken = (
            '$2b$12$wcwJ3EJ2zfq/FiEmy9YupejrQtsXIpqde5ZSzvuunF67eFvFuaKAu'
        )

    def test_success(self):
        request_data = {'username': self.user.username,
                        'token': self.validToken}
        response = self.client.post(self.url, request_data)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data, {
            'msg': 'Token verification successful.',
            'username': self.user.username
        })

