from rest_framework import serializers

from ..models import Board


class BoardSerializer(serializers.ModelSerializer):
    class Meta:
        model = Board
        fields = ('id', 'team', 'name')
        extra_kwargs = {
            'name': {
                'error_messages': {
                    'blank': 'Name cannot be blank.',
                    'null': 'Name cannot be null.'
                }
            },
            'team': {
                'error_messages': {
                    'blank': 'Team cannot be blank.',
                    'null': 'Team cannot be null.',
                    'does_not_exist': 'Team does not exist.'
                }
            }
        }

