<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::table('greenprints', function (Blueprint $table) {
            $table->foreignId("item_id")->references("id")->on("items");
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::table('greenprints', function (Blueprint $table) {
            $table->dropForeign(['greenprint_id']);
            $table->dropColumn('greenprint_id');
        });
    }
};
