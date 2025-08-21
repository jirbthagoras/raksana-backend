<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use Illuminate\Support\Facades\Hash;
use App\Models\User;

class MakeFilamentUser extends Command
{
    protected $signature = 'make:custom-filament-user';
    protected $description = 'Create a new Filament user with username instead of name';

    public function handle(): int
    {
        $username = $this->ask('Username');
        $email = $this->ask('Email address');
        $password = $this->secret('Password');

        User::create([
            'username' => $username,   // 👈 fixed
            'email' => $email,
            'password' => Hash::make($password),
        ]);

        $this->info("Filament user `{$username}` created successfully!");

        return self::SUCCESS;
    }
}
