<?php

namespace Tests\Feature\API\user;

use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Testing\Fluent\AssertableJson;
use Tests\TestCase;

class RegisterControllerTest extends TestCase
{
    use RefreshDatabase;

    public function test_register_endpoint()
    {
        $user = User::factory(1)->makeOne();

        $response = $this->postJson('/api/auth/register', [
            'name' => $user->name,
            'email' => $user->email,
            'password' => 'password',
        ]);

        $response->assertStatus(200);

        $response->assertJson(function (AssertableJson $json) use($user){

            $json->hasAll('data.user');

            $json->whereAll([
                'data.user.name' => $user['name'],
                'data.user.email' => $user['email']
            ])->etc();
        });
    }

    public function test_fail_register_endpoint()
    {

        $response = $this->postJson('/api/auth/register', [
            'name' => 'test name',
            'password' => 'password',
        ]);

        $response->assertStatus(401);

        $response->assertJson(function (AssertableJson $json) {

            $json->hasAll('error.message');

            $json->whereAll([
                'error.message' => 'Error to create a new user',
            ])->etc();
        });
    }
}
