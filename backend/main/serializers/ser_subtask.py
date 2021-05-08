from rest_framework import serializers
from ..models import Subtask
from ..validation.val_auth import authenticate_custom, authorize_custom, \
    authorization_error
from ..validation.val_custom import CustomAPIException
import status


class SubtaskSerializer(serializers.ModelSerializer):
    title = serializers.CharField(
        max_length=50,
        error_messages={
            'max_length':
                'Subtask titles cannot be longer than 50 characters.',
            'blank':
                'Subtask title cannot be empty.'
        }
    )

    class Meta:
        model = Subtask
        fields = '__all__'


# only to be used by PATCH actions
class SubtaskUpdateSerializer(SubtaskSerializer):
    @staticmethod
    def validate_id(value):
        if not value:
            raise CustomAPIException('id',
                                   'Subtask ID cannot be empty.',
                                     status.HTTP_400_BAD_REQUEST)
        try:
            subtask = Subtask.objects.select_related(
                'task',
                'task__user',
                'task__column__board'
            ).get(id=value)
        except Subtask.DoesNotExist:
            raise CustomAPIException('id',
                                   'Subtask not found.',
                                     status.HTTP_404_NOT_FOUND)
        return subtask

    @staticmethod
    def validate_data(value):
        if not value:
            raise CustomAPIException('data',
                                   'Data cannot be empty.',
                                     status.HTTP_400_BAD_REQUEST)
        return value

    @staticmethod
    def validate_title(value):
        if not value:
            raise CustomAPIException('title',
                                   'Title cannot be empty.',
                                     status.HTTP_400_BAD_REQUEST)

    @staticmethod
    def validate_done(value):
        if value == '' or value is None or not value:
            raise CustomAPIException('done',
                                   'Done cannot be empty.',
                                     status.HTTP_400_BAD_REQUEST)

    @staticmethod
    def validate_order(value):
        if value == '' or value is None or not value:
            raise CustomAPIException('order',
                                   'Order cannot be empty.',
                                     status.HTTP_400_BAD_REQUEST)

    def validate(self, attrs):
        auth_user = attrs.get('user')

        user, authentication_error = authenticate_custom(
            auth_user.get('username'),
            auth_user.get('token')
        )
        if authentication_error:
            raise authentication_error

        subtask = self.validate_id(attrs.get('id'))

        authorization_res = authorize_custom(auth_user.get('username'))
        if authorization_res and subtask.task.user != user \
                or subtask.task.column.board.team_id != user.team.id:
            raise authorization_error

        data = self.validate_data(attrs.get('data'))
        if 'title' in data.keys():
            self.validate_title(data.get('title'))
        if 'done' in data.keys():
            self.validate_done(data.get('done'))
        if 'order' in data.keys():
            self.validate_order(data.get('order'))

        return {'instance': subtask, 'data': data}

    def update(self, instance, validated_data):
        serializer = SubtaskSerializer(instance,
                                       data=validated_data,
                                       partial=True)
        serializer.is_valid(raise_exception=True)
        serializer.save()
        return {'msg': 'Subtask update successful.',
                'id': serializer.data.get('id')}
