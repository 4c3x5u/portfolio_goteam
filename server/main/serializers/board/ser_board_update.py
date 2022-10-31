from rest_framework import serializers

from .ser_board import BoardSerializer
from server.main.helpers.auth_helper import AuthHelper
from ...models import Board


class UpdateBoardSerializer(serializers.ModelSerializer):
    board = serializers.PrimaryKeyRelatedField(
        queryset=Board.objects.all(),
        error_messages={'null': 'Board ID cannot be null.',
                        'incorrect_type': 'Board ID must be a number.',
                        'does_not_exist': 'Board does not exist.'}
    )
    payload = serializers.DictField(allow_empty=False)
    auth_user = serializers.CharField(allow_blank=True)
    auth_token = serializers.CharField(allow_blank=True)

    class Meta:
        model = BoardSerializer.Meta.model
        fields = 'board', 'payload', 'auth_user', 'auth_token'

    def validate(self, attrs):
        user = AuthHelper.authenticate(attrs.get('auth_user'),
                                       attrs.get('auth_token'))
        board = attrs.get('board')
        AuthHelper.authorize(user, board.team_id)

        payload = attrs.get('payload')
        board_serializer = BoardSerializer(board, data=payload, partial=True)
        board_serializer.is_valid(raise_exception=True)

        self.instance = board
        return payload

    def to_representation(self, instance):
        return {'msg': 'Board updated successfuly.',
                'id': instance.id}
