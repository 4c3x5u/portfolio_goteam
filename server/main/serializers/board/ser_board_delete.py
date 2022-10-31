from rest_framework import serializers
import status

from server.main.serializers.board.ser_board import BoardSerializer
from server.main.helpers.auth_helper import AuthHelper
from server.main.helpers.custom_api_exception import CustomAPIException
from server.main.models import Board


class DeleteBoardSerializer(serializers.ModelSerializer):
    board = serializers.PrimaryKeyRelatedField(
        queryset=Board.objects.prefetch_related('team__board_set').all(),
        error_messages={'null': 'Board ID cannot be null.',
                        'incorrect_type': 'Board ID must be a number.',
                        'does_not_exist': 'Board does not exist.'}
    )
    auth_user = serializers.CharField(allow_blank=True)
    auth_token = serializers.CharField(allow_blank=True)

    class Meta:
        model = BoardSerializer.Meta.model
        fields = 'board', 'auth_user', 'auth_token',

    def validate(self, attrs):
        user = AuthHelper.authenticate(attrs.get('auth_user'),
                                       attrs.get('auth_token'))
        board = attrs.get('board')
        if len(board.team.board_set.all()) <= 1:
            raise CustomAPIException(
                'board',
                'You cannot delete the last remaining board.',
                status.HTTP_400_BAD_REQUEST)
        AuthHelper.authorize(user, board.team_id)
        return board

    def delete(self):
        self.instance = {'id': self.validated_data.id}
        return self.validated_data.delete()

    def to_representation(self, instance):
        return {
            'msg': 'Board deleted successfully.',
            'id': instance.get('id'),
        }



