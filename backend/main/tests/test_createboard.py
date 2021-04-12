from rest_framework.test import APITestCase
from rest_framework.exceptions import ErrorDetail
from ..models import Board, Team, User, Column
from ..util import new_member, new_admin, forbidden_response


class CreateBoardTests(APITestCase):
    def setUp(self):
        self.url = '/boards/'
        self.team = Team.objects.create()
        self.admin = new_admin(self.team)
        self.member = new_member(self.team)
        
    def test_success(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url,
                                    {'team_id': self.team.id},
                                    HTTP_AUTH_USER=self.admin['username'],
                                    HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.data.get('msg'),
                         'Board creation successful.')
        board = Board.objects.get(id=response.data.get('board_id'))
        columns = Column.objects.filter(board=board.id)
        self.assertEqual(len(columns), 4)
        self.assertEqual(Board.objects.count(), initial_count + 1)

    def test_user_not_admin(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url,
                                    {'team_id': self.team.id},
                                    HTTP_AUTH_USER=self.member['username'],
                                    HTTP_AUTH_TOKEN=self.member['token'])
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, {
            'auth': ErrorDetail(string='The user is not an admin.',
                                code='not_authorized')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_team_id_blank(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url,
                                    {'team_id': ''},
                                    HTTP_AUTH_USER=self.admin['username'],
                                    HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team ID cannot be empty.',
                                   code='blank')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_team_not_found(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url,
                                    {'team_id': '123'},
                                    HTTP_AUTH_USER=self.admin['username'],
                                    HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 404)
        self.assertEqual(response.data, {
            'team_id': ErrorDetail(string='Team not found.', code='not_found')
        })
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_token_empty(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url,
                                    {'team_id': self.team.id},
                                    HTTP_AUTH_USER=self.admin['username'],
                                    HTTP_AUTH_TOKEN='')
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, forbidden_response)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_token_invalid(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url,
                                    {'team_id': self.team.id},
                                    HTTP_AUTH_USER=self.admin['username'],
                                    HTTP_AUTH_TOKEN='ASDKFJ!FJ_012rjpiwajfosi')
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, forbidden_response)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_user_blank(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url,
                                    {'team_id': self.team.id},
                                    HTTP_AUTH_USER='',
                                    HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, forbidden_response)
        self.assertEqual(Board.objects.count(), initial_count)

    def test_auth_user_invalid(self):
        initial_count = Board.objects.count()
        response = self.client.post(self.url,
                                    {'team_id': self.team.id},
                                    HTTP_AUTH_USER='invalidio',
                                    HTTP_AUTH_TOKEN=self.admin['token'])
        self.assertEqual(response.status_code, 403)
        self.assertEqual(response.data, forbidden_response)
        self.assertEqual(Board.objects.count(), initial_count)


