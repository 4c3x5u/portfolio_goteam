from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Board, Team, Column
from ..util import new_member, new_admin, not_authenticated_response_data


class CreateBoardTests(APITestCase):
    endpoint = '/boards/'

    def setUp(self):
        self.team = Team.objects.create()
        self.member = new_member(self.team)
        self.admin = new_admin(self.team)

    def test_success(self):
        initial_count = Board.objects.count()
        response = self.client.post(
            self.endpoint,
            {'team_id': self.team.id, 'name': 'My Board'},
            HTTP_AUTHORIZATION=self.admin['auth_header']
        )
        print(f'§response: {response.data}')
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.data.get('msg'),
                         'Board creation successful.')
        board = Board.objects.get(id=response.data.get('board_id'))
        columns = Column.objects.filter(board=board.id)
        self.assertEqual(len(columns), 4)
        self.assertEqual(Board.objects.count(), initial_count + 1)

    def test_user_not_admin(self):
        initial_count = Board.objects.count()
        response = self.client.post(
            self.endpoint,
            {'team_id': self.team.id},
            HTTP_AUTHORIZATION=self.member['auth_header']
        )
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, {
            'auth': ErrorDetail(string='The user is not an admin.',
                                code='not_authorized')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_team_id_blank(self):
        initial_count = Board.objects.count()
        response = self.client.post(
            self.endpoint,
            {'team_id': ''},
            HTTP_AUTHORIZATION=self.admin['auth_header']
        )
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team ID cannot be empty.',
                                   code='blank')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_team_not_found(self):
        initial_count = Board.objects.count()
        response = self.client.post(
            self.endpoint,
            {'team_id': '123'},
            HTTP_AUTHORIZATION=self.admin['auth_header']
        )
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team not found.', code='not_found')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_token_empty(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.endpoint,
                                    {'team_id': self.team.id},
                                    HTTP_AUTH_USER=self.admin['username'])
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response_data)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_token_invalid(self):
        initial_count = Board.objects.count()
        response = self.client.post(
            self.endpoint,
            {'team_id': self.team.id},
            HTTP_AUTHORIZATION=f'{self.admin["username"]} ASDf/lasdkfajsdflalx'
        )
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response_data)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_user_blank(self):
        initial_count = Board.objects.count()
        response = self.client.post(
            self.endpoint,
            {'team_id': self.team.id},
            HTTP_AUTH_TOKEN=self.admin['auth_header'].split()[1]
        )
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response_data)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_user_invalid(self):
        initial_count = Board.objects.count()
        response = self.client.post(
            self.endpoint,
            {'team_id': self.team.id},
            HTTP_AUTH_TOKEN=f'bbanovich {self.admin["auth_header"].split()[1]}'
        )
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, not_authenticated_response_data)
        self.assertEqual(Board.objects.count(), initial_count)


